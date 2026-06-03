# Git Push Timer 开发说明

本文面向开发者，记录当前有效的架构说明、模块职责、构建方式和关键设计决策。

过程性的排查、复盘、日志证据和验证记录，统一放在 `docs/investigations/` 下，不在本文中按时间线展开。

---

## 项目概述

`git-push-timer` 是一个跨平台（macOS + Windows）的本地目录定时 Git 同步工具，主要用于：
- 监控多个本地目录（如 Obsidian 笔记、Bruno 集合等）
- 自动检测变更并执行 `git commit` + `git push`
- 将本地数据通过 Git 仓库进行版本控制和远程备份
- 保持数据存储在用户自己的 Git 仓库中

---

## 技术选型

| 方案 | 选择 | 理由 |
|------|------|------|
| 语言 | Go 1.21 | 跨平台编译、单文件部署、适合做轻量命令行工具 |
| 定时调度 | robfig/cron/v3 | 用于解析 Cron 表达式并计算下一次执行时间，不依赖系统定时任务 |
| 日志 | 文件日志 | 输出到 `<可执行文件所在目录>/logs/` |

---

## 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                  Git Push Timer (Go)                    │
│  ┌─────────────────────────────────────────────────┐   │
│  │       Scheduler（Cron 解析 + 60s 轮询）          │   │
│  │    维护每个仓库的 nextRun，并按整分钟检查         │   │
│  └─────────────────────┬───────────────────────────┘   │
│                        │                                 │
│                        ▼                                 │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Executor                           │   │
│  │  1. 接收已到期仓库任务                           │   │
│  │  2. 调用 Git 命令执行                            │   │
│  │  3. 写入日志                                     │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

---

## 核心模块

### 1. 配置读取 (`internal/config/config.go`)
- 读取可执行文件同级目录下的 `config/repos.json`
- 支持多个仓库配置
- 支持启用/禁用开关

### 2. Git 执行 (`internal/executor/executor.go`)
- 检查目录是否存在
- 检查是否是 Git 仓库
- 使用 `git status --porcelain` 检测变更（包含未跟踪文件）
- `git add` + `git commit` + `git push`
- 空提交保护

### 3. 日志记录 (`internal/logger/logger.go`)
- 日志输出到 `<可执行文件目录>/logs/`
- 按日期命名日志文件
- 同时输出到控制台

### 4. 定时调度 (`internal/scheduler/scheduler.go`)
- 支持每个仓库独立配置标准 5 段 Cron 表达式
- 按整分钟进行 `60s` 轮询，检查是否有仓库到达 `nextRun`
- 启动后从下一个匹配的 Cron 时间点开始检查，不会立即执行一次
- 机器睡眠错过多个计划时间点时只补跑 1 次
- 同一仓库上一次任务未结束时跳过本次执行并记录日志
- 支持停止

---

## 项目结构

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
├── docs/
│   ├── investigations/                      # 专题排查、复盘与验证记录
│   └── issue-ledger.md                      # 跨专题问题台账
├── README.md                                # 项目说明
├── DEVELOPMENT.md                           # 开发说明
├── CLAUDE.md                                # Claude Code 协作说明
├── AGENTS.md                                # Codex 自定义指令与工作流规范
├── build.sh                                 # 本地构建脚本
├── release.sh                               # 发布打包脚本
├── LICENSE
├── go.mod
└── go.sum
```

---

## 构建命令

本地构建可直接运行 [build.sh](./build.sh)：

```bash
./build.sh
```

该脚本会下载依赖，并在项目根目录生成 macOS 与 Windows 可执行文件：
- `git-push-timer`
- `git-push-timer.exe`

发布打包使用 [release.sh](./release.sh)：

```bash
./release.sh v1.0.0
```

该脚本会将版本号写入 `main.version`，并在 `dist/` 目录生成带版本号的 zip 包，同时复制 `config/repos.json.example`。

也可以手动执行以下命令：

```bash
# 下载依赖
go mod download

# macOS 编译
GOOS=darwin GOARCH=amd64 go build -o git-push-timer ./cmd/git-push-timer

# Windows 编译
GOOS=windows GOARCH=amd64 go build -o git-push-timer.exe ./cmd/git-push-timer
```

---

## 开发注意事项

- 默认 Cron 表达式为 `*/5 * * * *`（每 5 分钟执行一次），定义在 `internal/scheduler/scheduler.go` 中的 `defaultCronSpec` 常量。
- 修改默认频率需要调整 `defaultCronSpec`；单个仓库的频率可通过配置文件中的 `cronSpec` 字段覆盖。
- 当前配置仅在调度器启动时读取；修改 `enabled` 或 `cronSpec` 后需要重启程序才能生效。
- `config/repos.json` 与 `logs/` 属于本地运行文件，已在 `.gitignore` 中，不会被提交到 Git。

---

## 相关链接

- GitHub 仓库：https://github.com/Miracle-ZT/git-push-timer
- 用户文档：[README.md](./README.md)
- Claude Code 记忆：[CLAUDE.md](./CLAUDE.md)
- Codex 指令：[AGENTS.md](./AGENTS.md)
- 问题台账：[docs/issue-ledger.md](./docs/issue-ledger.md)
