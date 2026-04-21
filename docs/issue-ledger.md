# Issue Ledger

说明：本文档用于长期维护跨专题的问题台账。专项排查过程、日志证据和完整时间线仍放在 `docs/investigations/` 下；如果某个问题没有单独的排查文档，也可以直接在这里登记和更新状态。

## 状态约定

- `已解决`：代码和相关文档已完成修复，并已落到仓库提交中。
- `可选优化`：当前实现可接受，但仍有进一步优化空间。

## 已解决

| 问题 | 状态 | 关联文档 | 关联提交 | 备注 |
| --- | --- | --- | --- | --- |
| macOS 睡眠后长周期 Cron 调度严重漂移 | 已解决 | `docs/investigations/2026-04-16-macos-sleep-scheduler.md` | `2aa02dc` | 调度改为按整分钟 `60s` 轮询；错过多个时间点时只补跑 1 次 |
| 日志未按自然日轮转 | 已解决 | `docs/investigations/2026-04-16-macos-sleep-scheduler.md` | `7dd87ec` | 按写入时的自然日切换日志文件 |
| 未跟踪文件不会被识别为 Git 变更 | 已解决 | `docs/investigations/2026-04-16-macos-sleep-scheduler.md` | `29b8c24` | 变更检测改为基于 `git status --porcelain` |
| Git 状态判断把“命令失败”和“有变更”混在一起 | 已解决 | `docs/investigations/2026-04-16-macos-sleep-scheduler.md` | `29b8c24` | 统一按 porcelain 输出判断仓库状态 |
| `os.Chdir` 带来的并发工作目录竞争 | 已解决 | `docs/investigations/2026-04-16-macos-sleep-scheduler.md` | `89ff4d0` | 每次执行命令时单独设置 `cmd.Dir`，不再修改进程全局目录 |
| 退出时未等待运行中的任务安全收尾 | 已解决 | `docs/investigations/2026-04-16-macos-sleep-scheduler.md` | `a9a1ee8` | 调度器停止时等待正在执行的任务结束后再退出 |
| “启动后立即执行一次检查”的文档与实现不一致 | 已解决 | `docs/investigations/2026-04-16-macos-sleep-scheduler.md` | `2aa02dc`, `555f046` | README 和开发文档已改为“从下一个匹配时间点开始检查” |
| `.gitignore` 中的构建产物规则误匹配源码目录 | 已解决 | 无单独专题文档 | `3f8fe6c` | 将忽略规则收紧为仓库根目录下的构建产物文件 |

## 可选优化

| 问题 | 状态 | 关联文档 | 备注 |
| --- | --- | --- | --- |
| 系统从睡眠唤醒后立即主动触发一次调度检查 | 可选优化 | `docs/investigations/2026-04-16-macos-sleep-scheduler.md` | 当前方案已解决长周期严重漂移，但第一次补跑仍可能晚几十秒 |
