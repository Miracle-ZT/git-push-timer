# AGENTS.md

## 项目概述

**Git Push Timer** — 跨平台（macOS + Windows）本地目录定时 Git 同步工具。

核心功能：监控用户指定的本地目录，自动检测变更并执行 `git commit` + `git push`，实现数据的版本控制和远程备份。

## 技术栈

- **语言**: Go 1.21
- **依赖**: `github.com/robfig/cron/v3`（Cron 表达式解析与下一次执行时间计算）

## 常用命令

```bash
# 本地构建 macOS 和 Windows 可执行文件
./build.sh

# 发布打包，生成 dist/ 下的 zip 包
./release.sh v1.0.0
```

手动构建命令、版本注入和打包产物说明见 [DEVELOPMENT.md](./DEVELOPMENT.md)。

## 开发约束

- 保持现有模块边界：`config` 负责配置读取与解析，`scheduler` 负责定时调度，`executor` 负责 Git 检查、提交和推送，`logger` 负责日志输出。
- 每个仓库可独立配置标准 5 段 Cron 表达式，默认值为 `*/5 * * * *`。
- 调度器按整分钟进行 `60s` 轮询，启动后从下一个匹配的 Cron 时间点开始检查，不应立即执行一次检查。
- 机器睡眠或程序暂停导致错过多个计划时间点时，同一仓库只补跑 1 次。
- 同一仓库上一次任务未结束时，应跳过本次执行并记录日志。
- 当前配置仅在调度器启动时读取；不要在未明确需求时引入配置热加载行为。

## 文档分工

- `README.md` 面向使用者，主要记录安装、配置、运行方式和用户可见行为。
- `DEVELOPMENT.md` 面向开发者，主要记录当前有效的架构说明、模块职责、构建方式和关键设计决策，不作为时间顺序的开发流水账。
- `docs/investigations/` 用于保存围绕单一主题展开的问题排查、复盘、review、验证过程、日志证据和修复思路等专题文档。
- `docs/issue-ledger.md` 用于长期维护跨专题的问题台账，记录问题状态、关联文档、关联提交和可选优化项。
- `AGENTS.md` 用于给 Codex 提供仓库级指令、工作流规范和文档约束，也是 Claude Code 通过 `CLAUDE.md` 引用的共享指令来源。
- `CLAUDE.md` 用于 Claude Code 加载仓库指令，目前通过 `@AGENTS.md` 引用本文件，不单独维护重复内容。

## 文档脱敏规则

- 以后凡是需要提交到仓库的排查记录、复盘文档、review 记录、问题分析文档，默认先做脱敏再提交。
- 需要脱敏的内容包括：本机绝对路径、用户名、进程号、设备/机器标识、个人目录信息，以及没有必要公开的本地仓库名称或路径标识。
- 优先保留技术语义，使用稳定占位符替换敏感内容，例如：`<dev-repo-root>`、`<release-root>`、`<repo-name>`、`<target-repo-path>`、`<observed-pid>`。
- 相对源码路径、时间线、日志模式、错误信息、结论和修复思路通常应保留，因为这些信息对复盘有价值。
- 如需保留未脱敏原始版本，只保留在仓库外的本地位置，不提交到 Git，不放在仓库内等待 `.gitignore` 保护。

## Commit Message 规则

- commit message 必须简洁、清晰，聚焦于本次提交的核心变化。
- 提交标题必须使用统一格式：`<type>(<scope>): <description>`
- `scope` 为必填项，必须明确指出本次提交影响的模块、目录或功能域。
- `scope` 应优先使用稳定、具体的名称，例如：`scheduler`、`executor`、`config`、`logger`、`readme`、`release`。
- 避免使用过于宽泛、临时或含义不清的 `scope`，例如：`misc`、`update`、`优化`、`调整`。
- `type` 推荐使用：`feat`、`fix`、`refactor`、`docs`、`test`、`build`、`ci`、`chore`。
- `description` 只概括“这次改了什么、解决了什么”，避免堆叠实现细节或把实现过程写进标题。
- 多模块同时改动时，优先填写本次提交的主要影响范围，不必为附带修改罗列多个 `scope`。
- 正文（Body）仅在确有必要时使用，简洁说明为什么要做这个改动，而不是详细描述如何实现。
- 正文（Body）使用列表时，列表项之间不要插入空行；命令行提交时应将连续列表放在同一个正文参数中，避免多个 `-m` 自动拆成多个段落。
- 避免过于冗长的标题和正文，优先让团队成员快速理解本次提交的目的和结果。

### 提交格式

仅标题：

```text
<type>(<scope>): <description>
```

包含正文：

```text
<type>(<scope>): <description>

- <body item>
- <body item>
```

### 示例

- `feat(scheduler): 新增按 cron 时间点触发检查`
- `fix(executor): 修复未跟踪文件未被提交的问题`
- `docs(readme): 更新定时检查行为说明`
- `refactor(config): 简化配置加载流程`
