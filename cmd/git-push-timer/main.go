package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git-push-timer/internal/config"
	"git-push-timer/internal/executor"
	"git-push-timer/internal/logger"
	"git-push-timer/internal/scheduler"
)

// 版本号，通过 -ldflags 注入
var version = "dev"

// 支持 --version 参数
func init() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("git-push-timer %s\n", version)
		os.Exit(0)
	}
}

func main() {
	// 初始化日志
	log, err := logger.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败：%v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Info("Git Push Timer %s", version)

	// 加载配置（初次启动时加载）
	cfg, err := config.Load()
	if err != nil {
		log.Error("加载配置失败：%v", err)
		log.Info("提示：请确保 config/repos.json 配置文件存在")
		os.Exit(1)
	}
	log.Info("已加载 %d 个仓库配置", len(cfg.Repositories))

	// 创建执行器和调度器
	exec := executor.New(log)
	sched := scheduler.New(log, exec)

	// 启动调度器（为每个仓库创建独立的定时任务）
	if err := sched.Start(); err != nil {
		log.Error("启动调度器失败：%v", err)
		os.Exit(1)
	}

	// 等待退出信号
	log.Info("按 Ctrl+C 退出程序")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("正在退出...")
	sched.Stop()
	time.Sleep(500 * time.Millisecond)
	log.Info("Git Push Timer 已退出")
}
