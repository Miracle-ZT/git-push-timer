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
- macOS: `git-push-timer_darwin_amd64.zip`
- Windows: `git-push-timer_windows_amd64.zip`

下载后解压，得到可执行文件 `git-push-timer` 或 `git-push-timer.exe`。

**查看版本：**

```bash
./git-push-timer --version
```

### 2. 配置

在可执行文件同级目录下创建 `config/repos.json`：

```json
{
  "repositories": [
    {
      "name": "Obsidian",
      "path": "~/ObsidianVault",
      "branch": "main",
      "enabled": true,
      "cronSpec": "*/5 * * * *"
    },
    {
      "name": "Bruno",
      "path": "/Users/yourname/Documents/BrunoCollections",
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
| `path` | 本地目录路径（支持 `~` 和绝对路径） |
| `branch` | 推送的分支名 |
| `enabled` | 是否启用 |
| `cronSpec` | 可选，标准 5 段 Cron 表达式，控制该仓库的检查频率。不配置则使用默认值（每 5 分钟） |

**Cron 表达式示例：**
| 表达式 | 含义 |
|--------|------|
| `*/5 * * * *` | 每 5 分钟 |
| `0 * * * *` | 每小时整点 |
| `0 */2 * * *` | 每 2 小时 |
| `0 9 * * *` | 每天早上 9 点 |
| `0 */6 * * *` | 每 6 小时 |

**注意：**
- 配置文件修改后，需要重启程序才能生效（包括 `enabled` 和 `cronSpec`）
- 仅支持标准 5 段 Cron 表达式：`分 时 日 月 周`
- 不支持 `@every 30s`、`@daily` 这类 descriptor 语法

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
1. 从下一个匹配的 Cron 时间点开始检查
2. 之后按 Cron 配置的时间点自动检查
3. 按 Ctrl+C 退出

### 5. 日志

日志文件位于 `logs/` 目录下（可执行文件同级目录），按日期命名：
```
logs/
  2026-03-31.log
  2026-04-01.log
```

`logs/` 目录已在 `.gitignore` 中，不会被提交到 Git。

## 自行编译

```bash
# 安装依赖
go mod download

# macOS 编译
GOOS=darwin GOARCH=amd64 go build -o git-push-timer ./cmd/git-push-timer

# Windows 编译
GOOS=windows GOARCH=amd64 go build -o git-push-timer.exe ./cmd/git-push-timer
```

编译后的二进制文件是独立的，无需安装 Go 运行时即可运行。

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
├── .gitignore               # Git 忽略文件
├── go.mod
├── go.sum
└── README.md
```

## 开发注意事项

- 默认 Cron 表达式：`*/5 * * * *`（每 5 分钟执行一次），定义在 `internal/scheduler/scheduler.go` 中的 `defaultCronSpec` 常量
- 修改频率需要改代码：调整 `defaultCronSpec` 的值，或修改配置文件中的 `cronSpec` 字段
- 配置生效：修改 `enabled` 或 `cronSpec` 后需要重启程序
- `config/repos.json` 是本地配置文件，已在 `.gitignore` 中，不会被提交
