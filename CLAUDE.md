# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Git Push Timer** — 跨平台（macOS + Windows）本地目录定时 Git 同步工具。

核心功能：监控用户指定的本地目录，自动执行 Git commit 和 push，实现数据的版本控制和远程备份。

## Tech Stack

- **语言**: Go 1.21
- **依赖**: `github.com/robfig/cron/v3`（Cron 表达式解析与下一次执行时间计算）

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
│   └── git-push-timer/                       # 程序入口
├── internal/
│   ├── config/                              # 配置读取与解析
│   ├── executor/                            # Git 检查、提交、推送
│   ├── logger/                              # 日志输出
│   └── scheduler/                           # 定时调度
├── config/
│   └── repos.json.example                   # 配置示例
├── docs/                                    # 排查记录与补充文档
├── README.md                                # 项目说明
├── DEVELOPMENT.md                           # 开发说明
├── CLAUDE.md                                # Claude Code 协作说明
├── AGENTS.md                                # Codex 自定义指令与工作流规范
├── build.sh                                 # 本地构建脚本
├── release.sh                               # 发布打包脚本
├── .gitignore                               # Git 忽略文件
├── go.mod
└── go.sum
```

## Architecture

1. **main.go** 启动后创建日志、执行器、调度器
2. **scheduler** 读取仓库配置，解析标准 5 段 Cron 表达式，维护每个仓库的 `nextRun`
3. **scheduler** 按整分钟执行 `60s` 轮询，检查是否有仓库到期，并在需要时补跑 1 次
4. **executor** 在仓库到期时执行 Git 检查、提交和推送
5. **logger** 将日志输出到 `<可执行文件目录>/logs/`

## Key Design Decisions

- 日志路径：可执行文件同级目录下的 `logs/` 子目录
- 定时频率：支持每个仓库独立配置标准 5 段 Cron 表达式，默认为 `*/5 * * * *`
- 调度策略：使用 `60s` 轮询对齐整分钟，不在启动时立即执行一次检查
- 配置生效：修改 `enabled` 或 `cronSpec` 后需要重启程序
- 路径支持：`~` 开头路径会自动展开为用户主目录，也支持绝对路径
- `config/repos.json` 是本地配置文件，不提交到 Git
