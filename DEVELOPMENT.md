# Git Push Timer 开发记录

## 项目背景

用户需求：开发一个跨平台（macOS + Windows）的本地目录定时 Git 同步工具，用于：
- 监控多个本地目录（如 Obsidian 笔记、Bruno 集合等）
- 自动检测变更并执行 `git commit` + `git push`
- 实现数据的版本控制和远程备份
- 强调安全性和隐私性，数据存储在用户自己的 Git 仓库中

---

## 技术选型

| 方案 | 选择 | 理由 |
|------|------|------|
| 语言 | Go 1.21 | 用户有 Java 背景，学 Go 很快；跨平台编译；单文件部署 |
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

## 核心功能实现

### 1. 配置读取 (`internal/config/config.go`)
- 读取可执行文件同级目录下的 `config/repos.json`
- 支持多个仓库配置
- 支持启用/禁用开关

### 2. Git 执行 (`internal/executor/executor.go`)
- 检查目录是否存在
- 检查是否是 Git 仓库
- `git diff --quiet` 检测变更
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

## 问题与解决方案

### 问题：多仓库独立 PAT 配置

#### 背景
用户为每个 GitHub 仓库分别生成了独立的 PAT（Personal Access Token），每个 PAT 只允许访问对应的仓库（最小权限原则）。需要将多个 PAT 都存储到 macOS Keychain 中。

#### 问题
Git credential 存储基于 `protocol + host + username`，同一账号的多个仓库会共用一个 credential 条目，后输入的 PAT 会覆盖之前的。

#### 解决方案
在 remote URL 中使用不同的"用户名标记"来区分：

```bash
# 仓库 1
git remote set-url origin https://Miracle-ZT-git-push-timer@github.com/Miracle-ZT/git-push-timer.git
git push
# 输入 repo1 对应的 PAT

# 仓库 2
git remote set-url origin https://Miracle-ZT-other-repo@github.com/Miracle-ZT/other-repo.git
git push
# 输入 repo2 对应的 PAT
```

#### 原理
- Git credential 的"键" = `protocol + host + username`
- username 不同，credential 条目就不同
- GitHub 认证时只看 PAT 是否有效，不验证 username
- 钥匙串中会存储为独立条目：
  - `github.com - Miracle-ZT-git-push-timer`
  - `github.com - Miracle-ZT-other-repo`

#### 验证
打开 macOS **钥匙串访问** App，搜索 `github`，可以看到独立的条目。

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
├── docs/                                    # 排查记录与补充文档
├── README.md                                # 项目说明
├── DEVELOPMENT.md                           # 开发说明
├── CLAUDE.md                                # Claude Code 协作说明
├── AGENTS.md                                # Codex 自定义指令与工作流规范
├── build.sh                                 # 本地构建脚本
├── release.sh                               # 发布打包脚本
├── go.mod
└── go.sum
```

---

## 构建命令

```bash
# 下载依赖
go mod download

# macOS 编译
GOOS=darwin GOARCH=amd64 go build -o git-push-timer ./cmd/git-push-timer

# Windows 编译
GOOS=windows GOARCH=amd64 go build -o git-push-timer.exe ./cmd/git-push-timer
```

---

## 待办事项

- [ ] Windows 平台测试
- [ ] 路径 `~` 展开支持
- [ ] 失败通知功能
- [ ] CLI 管理命令（add/list/remove）
- [ ] GUI 界面（可选）

---

## 相关链接

- GitHub 仓库：https://github.com/Miracle-ZT/git-push-timer
- 用户文档：README.md
- 开发指南：CLAUDE.md
- Codex 指令：AGENTS.md
