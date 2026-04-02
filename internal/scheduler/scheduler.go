package scheduler

import (
	"git-push-timer/internal/config"
	"git-push-timer/internal/executor"
	"git-push-timer/internal/logger"

	"github.com/robfig/cron/v3"
)

// 默认 Cron 表达式（每 5 分钟）
const defaultCronSpec = "*/5 * * * *"

// Scheduler 定时调度器
type Scheduler struct {
	cron     *cron.Cron
	logger   *logger.Logger
	executor *executor.Executor
}

// New 创建调度器
func New(logger *logger.Logger, executor *executor.Executor) *Scheduler {
	return &Scheduler{
		cron:     cron.New(),
		logger:   logger,
		executor: executor,
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

	// 为每个启用的仓库创建独立的定时任务
	for _, repo := range cfg.Repositories {
		if !repo.Enabled {
			continue
		}

		// 使用仓库自己的 CronSpec，如果没有则用默认值
		cronSpec := repo.CronSpec
		if cronSpec == "" {
			cronSpec = defaultCronSpec
		}

		// 为这个仓库创建定时任务（捕获循环变量）
		repo := repo
		_, err := s.cron.AddFunc(cronSpec, func() {
			s.logger.Info("触发仓库：%s", repo.Name)
			if err := s.executor.ExecuteRepository(repo); err != nil {
				s.logger.Error("仓库 %s 处理失败：%v", repo.Name, err)
			}
		})
		if err != nil {
			s.logger.Error("为仓库 %s 创建定时任务失败：%v", repo.Name, err)
			continue
		}
		s.logger.Info("仓库 %s 已调度，频率：%s", repo.Name, cronSpec)
	}

	s.cron.Start()
	s.logger.Info("调度器已启动")
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.logger.Info("调度器已停止")
}
