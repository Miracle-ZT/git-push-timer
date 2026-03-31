# Git Push Timer

跨平台（macOS + Windows）本地目录定时 Git 同步工具。

## 功能

- 监控多个本地目录
- 自动检测变更并执行 `git commit` + `git push`
- 定时检查（默认每 5 分钟）
- 空提交保护（无变更时不执行）
- 文件日志记录

## 快速开始

### 1. 下载

从 [Releases](https://github.com/Miracle-ZT/git-push-timer/releases) 下载对应平台的二进制文件：
- macOS: `git-push-timer-darwin-amd64`
- Windows: `git-push-timer-windows-amd64.exe`

### 2. 配置

在可执行文件同级目录下创建 `config/repos.json`：

```json
{
  "repositories": [
    {
      "name": "Obsidian",
      "path": "~/ObsidianVault",
      "branch": "main",
      "enabled": true
    },
    {
      "name": "Bruno",
      "path": "~/Documents/BrunoCollections",
      "branch": "main",
      "enabled": true
    }
  ]
}
```

**配置说明：**
| 字段 | 说明 |
|------|------|
| `name` | 仓库名称（用于日志标识） |
| `path` | 本地目录路径（支持 `~` 简写） |
| `branch` | 推送的分支名 |
| `enabled` | 是否启用 |

### 3. 前置要求

- 目录必须已经初始化为 Git 仓库
- 已配置远程仓库（`git remote add origin ...`）
- 已配置 Git 认证（SSH key 或 Credential）

### 4. 运行

```bash
# macOS/Linux
./git-push-timer

# Windows
git-push-timer-windows-amd64.exe
```

程序启动后：
1. 立即执行一次检查
2. 之后每 5 分钟自动检查一次
3. 按 Ctrl+C 退出

### 5. 日志

日志文件位于 `logs/` 目录下（可执行文件同级目录），按日期命名：
```
logs/
  2026-03-31.log
  2026-04-01.log
```

## 自行编译

```bash
# 安装依赖
go mod download

# macOS 编译
GOOS=darwin GOARCH=amd64 go build -o git-push-timer ./cmd/git-push-timer

# Windows 编译
GOOS=windows GOARCH=amd64 go build -o git-push-timer.exe ./cmd/git-push-timer
```

## 项目结构

```
git-push-timer/
├── cmd/
│   └── git-push-timer/
│       └── main.go          # 入口
├── internal/
│   ├── config/              # 配置读取
│   ├── executor/            # Git 执行
│   ├── logger/              # 日志记录
│   └── scheduler/           # 定时调度
├── config/
│   └── repos.json.example   # 配置示例
└── go.mod
```

## 开发注意事项

- 默认 Cron 表达式：`*/5 * * * *`（每 5 分钟执行一次）
- 修改频率需要改代码：`cmd/git-push-timer/main.go` 中的 `cronSpec` 变量
