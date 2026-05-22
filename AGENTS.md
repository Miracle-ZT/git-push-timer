# AGENTS.md

## 文档分工

- `README.md` 面向使用者，主要记录安装、配置、运行方式和用户可见行为。
- `DEVELOPMENT.md` 面向开发者，主要记录当前有效的架构说明、模块职责、构建方式和关键设计决策，不作为时间顺序的开发流水账。
- `docs/investigations/` 用于保存围绕单一主题展开的问题排查、复盘、review、验证过程、日志证据和修复思路等专题文档。
- `docs/issue-ledger.md` 用于长期维护跨专题的问题台账，记录问题状态、关联文档、关联提交和可选优化项。
- `CLAUDE.md` 用于给 Claude Code 提供仓库上下文和开发约定，文档分工需与本文件保持一致。
- `AGENTS.md` 用于给 Codex 提供仓库级指令、工作流规范和文档约束。
- `AGENTS.md` 与 `CLAUDE.md` 中的共享约束需要同步维护；任一方新增或修改后，另一方也要同步更新。

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
