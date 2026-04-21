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

## Documentation Responsibilities

- `README.md` 面向使用者，记录安装、配置、运行方式和用户可见行为。
- `DEVELOPMENT.md` 面向开发者，记录当前有效的架构说明、模块职责、构建方式和关键设计决策，不作为开发流水账。
- `docs/investigations/` 用于保存排查过程、复盘记录、review 结论、日志证据和验证过程。
- `AGENTS.md` 是 Codex 的仓库级指令文件；处理文档分工和脱敏时应与其中规则保持一致。
- `AGENTS.md` 与 `CLAUDE.md` 中的共享约束需要同步维护；任一方新增或修改后，另一方也要同步更新。

## Document Redaction Rules

- 以后凡是需要提交到仓库的排查记录、复盘文档、review 记录、问题分析文档，默认先做脱敏再提交。
- 需要脱敏的内容包括：本机绝对路径、用户名、进程号、设备/机器标识、个人目录信息，以及没有必要公开的本地仓库名称或路径标识。
- 优先保留技术语义，使用稳定占位符替换敏感内容，例如：`<dev-repo-root>`、`<release-root>`、`<repo-name>`、`<target-repo-path>`、`<observed-pid>`。
- 相对源码路径、时间线、日志模式、错误信息、结论和修复思路通常应保留，因为这些信息对复盘有价值。
- 如需保留未脱敏原始版本，只保留在仓库外的本地位置，不提交到 Git，不放在仓库内等待 `.gitignore` 保护。

## Commit Message Rules

- commit message 要更加简洁、清晰，并聚焦于本次提交的核心变化。
- 避免过于冗长的描述，保留问题本身的核心和关键改动即可。
- 可以省略过于具体的实现细节，优先让团队成员快速理解“这次改了什么、解决了什么”。
- 保持结构清晰，正文只补充必要信息，不把实现过程写成过长的说明。
