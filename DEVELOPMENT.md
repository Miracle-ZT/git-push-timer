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
| 定时调度 | robfig/cron/v3 | 内置 cron，不依赖系统定时任务，跨平台一致 |
| 日志 | 文件日志 | 输出到 `<可执行文件所在目录>/logs/` |

---

## 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                  Git Push Timer (Go)                    │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Scheduler (cron 库)                 │   │
│  │           每 5 分钟触发一次执行                   │   │
│  └─────────────────────┬───────────────────────────┘   │
│                        │                                 │
│                        ▼                                 │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Executor                           │   │
│  │  1. 读取 config/repos.json                      │   │
│  │  2. 遍历所有仓库                                 │   │
│  │  3. 调用 Git 命令执行                            │   │
│  │  4. 写入日志                                     │   │
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
- 使用 cron 表达式 `*/5 * * * *`（每 5 分钟）
- 启动时立即执行一次
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
│   └── git-push-timer/
│       └── main.go          # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go        # 配置读取（repos.json）
│   ├── executor/
│   │   └── executor.go      # Git 执行逻辑（add/commit/push）
│   ├── logger/
│   │   └── logger.go        # 日志记录（输出到 logs/目录）
│   └── scheduler/
│       └── scheduler.go     # 定时调度（cron）
├── config/
│   └── repos.json.example   # 配置文件示例
├── go.mod
├── go.sum
├── README.md
└── CLAUDE.md
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
