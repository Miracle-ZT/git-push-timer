# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Git Push Timer** — 跨平台（macOS + Windows）本地目录定时 Git 同步工具。

核心功能：监控用户指定的本地目录，自动执行 Git commit 和 push，实现数据的版本控制和远程备份。

## Tech Stack

- **语言**: Go 1.21
- **依赖**: `github.com/robfig/cron/v3`（定时调度）

## Build Commands

```bash
# 下载依赖
go mod download

# macOS 编译
GOOS=darwin GOARCH=amd64 go build -o git-push-timer ./cmd/git-push-timer

# Windows 编译
GOOS=windows GOARCH=amd64 go build -o git-push-timer.exe ./cmd/git-push-timer
```

## Project Structure

```
git-push-timer/
├── cmd/
│   └── git-push-timer/
│       └── main.go          # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go        # 配置读取（repos.json），支持 ~ 路径展开
│   ├── executor/
│   │   └── executor.go      # Git 执行逻辑（add/commit/push）
│   ├── logger/
│   │   └── logger.go        # 日志记录（输出到 logs/目录）
│   └── scheduler/
│       └── scheduler.go     # 定时调度（cron），为每个仓库创建独立任务
├── config/
│   └── repos.json.example   # 配置文件示例
├── .gitignore               # Git 忽略文件
├── go.mod
├── go.sum
└── README.md
```

## Architecture

1. **main.go** 启动后加载配置，创建日志、执行器、调度器
2. **scheduler** 为每个仓库创建独立的 cron 定时任务
3. **executor** 遍历所有配置的仓库，检测变更并执行 git push
4. **logger** 将日志输出到 `<可执行文件目录>/logs/`

## Key Design Decisions

- 日志路径：可执行文件同级目录下的 `logs/` 子目录
- 定时频率：支持每个仓库独立配置 Cron 表达式，默认为 `*/5 * * * *`（每 5 分钟）
- 配置生效：修改 `enabled` 或 `cronSpec` 后需要重启程序
- 路径支持：`~` 开头路径会自动展开为用户主目录，也支持绝对路径
- `config/repos.json` 是本地配置文件，不提交到 Git
