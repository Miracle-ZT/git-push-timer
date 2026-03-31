package scheduler

import (
	"git-push-timer/internal/config"
	"git-push-timer/internal/executor"
	"git-push-timer/internal/logger"

	"github.com/robfig/cron/v3"
)

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
func (s *Scheduler) Start(cfg *config.Config, spec string) error {
	// 解析 cron 表达式
	_, err := s.cron.AddFunc(spec, func() {
		s.RunAllRepositories(cfg)
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	s.logger.Info("调度器已启动，定时：%s", spec)
	return nil
}

// RunAllRepositories 遍历所有启用的仓库（公开方法）
func (s *Scheduler) RunAllRepositories(cfg *config.Config) {
	s.logger.Info("=== 开始检查所有仓库 ===")

	for _, repo := range cfg.Repositories {
		if !repo.Enabled {
			s.logger.Info("仓库 %s 已禁用，跳过", repo.Name)
			continue
		}

		if err := s.executor.ExecuteRepository(repo); err != nil {
			s.logger.Error("仓库 %s 处理失败：%v", repo.Name, err)
		}
	}

	s.logger.Info("=== 所有仓库检查完成 ===")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.logger.Info("调度器已停止")
}
