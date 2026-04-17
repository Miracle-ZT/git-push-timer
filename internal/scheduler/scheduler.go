package scheduler

import (
	"sync"
	"time"

	"git-push-timer/internal/config"
	"git-push-timer/internal/executor"
	"git-push-timer/internal/logger"

	"github.com/robfig/cron/v3"
)

// 默认 Cron 表达式（每 5 分钟）
const defaultCronSpec = "*/5 * * * *"

const pollInterval = time.Minute

var standardCronParser = cron.NewParser(
	cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

type scheduledRepo struct {
	repo     config.Repository
	schedule cron.Schedule
	nextRun  time.Time
	running  bool
	mu       sync.Mutex
}

// Scheduler 定时调度器
type Scheduler struct {
	logger   *logger.Logger
	executor *executor.Executor
	jobs     []*scheduledRepo
	stopCh   chan struct{}
	doneCh   chan struct{}
	stopOnce sync.Once
}

// New 创建调度器
func New(logger *logger.Logger, executor *executor.Executor) *Scheduler {
	return &Scheduler{
		logger:   logger,
		executor: executor,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		s.logger.Error("加载配置失败：%v", err)
		return err
	}

	now := time.Now()

	// 为每个启用的仓库创建调度状态
	for _, repo := range cfg.Repositories {
		if !repo.Enabled {
			continue
		}

		// 使用仓库自己的 CronSpec，如果没有则用默认值
		cronSpec := repo.CronSpec
		if cronSpec == "" {
			cronSpec = defaultCronSpec
		}

		schedule, err := parseCronSpec(cronSpec)
		if err != nil {
			s.logger.Error("为仓库 %s 创建定时任务失败：%v", repo.Name, err)
			continue
		}

		s.jobs = append(s.jobs, &scheduledRepo{
			repo:     repo,
			schedule: schedule,
			nextRun:  schedule.Next(now),
		})
		s.logger.Info("仓库 %s 已调度，频率：%s", repo.Name, cronSpec)
	}

	go s.run()

	s.logger.Info("调度器已启动")
	return nil
}

func (s *Scheduler) run() {
	defer close(s.doneCh)

	timer := time.NewTimer(timeUntilNextTick(time.Now()))
	defer timer.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-timer.C:
			s.checkDueRepositories()
			timer.Reset(timeUntilNextTick(time.Now()))
		}
	}
}

func (s *Scheduler) checkDueRepositories() {
	now := time.Now()

	for _, job := range s.jobs {
		job.mu.Lock()
		if now.Before(job.nextRun) {
			job.mu.Unlock()
			continue
		}

		job.nextRun = advanceNextRun(job.schedule, job.nextRun, now)
		if job.running {
			repoName := job.repo.Name
			job.mu.Unlock()
			s.logger.Warn("仓库 %s 已到执行时间，但上一次任务仍未结束，跳过本次", repoName)
			continue
		}

		job.running = true
		repo := job.repo
		job.mu.Unlock()

		s.logger.Info("触发仓库：%s", repo.Name)
		go s.execute(job, repo)
	}
}

func (s *Scheduler) execute(job *scheduledRepo, repo config.Repository) {
	defer func() {
		job.mu.Lock()
		job.running = false
		job.mu.Unlock()
	}()

	if err := s.executor.ExecuteRepository(repo); err != nil {
		s.logger.Error("仓库 %s 处理失败：%v", repo.Name, err)
	}
}

func parseCronSpec(cronSpec string) (cron.Schedule, error) {
	return standardCronParser.Parse(cronSpec)
}

func advanceNextRun(schedule cron.Schedule, nextRun, now time.Time) time.Time {
	for !nextRun.After(now) {
		nextRun = schedule.Next(nextRun)
	}
	return nextRun
}

func timeUntilNextTick(now time.Time) time.Duration {
	nextTick := now.Truncate(pollInterval).Add(pollInterval)
	wait := nextTick.Sub(now)
	if wait <= 0 {
		return pollInterval
	}
	return wait
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
		<-s.doneCh
		s.logger.Info("调度器已停止")
	})
}
