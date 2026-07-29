# 基于 module 依赖安装 skills

通过 Go module 依赖, 将 Agent Skills 安装到项目本地的 `.agents/skills` 目录.

实现与接入约定见 `github.com/xoctopus/genx/pkg/agent` (go doc).

## 接入要点

按照 `github.com/xoctopus/genx/pkg/agent` 约定接入:

1. 在直接 `require` 上方标注 `// +skill:<name>` (可多个)
2. 技能提供方在 `.agents/skills/<name>/` 提供内容 (如 `SKILL.md`)
3. 调用 `(&agent.Installer{}).Install(ctx)` 完成安装

## 约定补充

- 仅识别直接依赖, 忽略 `indirect`
- 安装结果为 symlink: `<project>/.agents/skills/<name>` -> 模块内同路径
- `.gitignore` 追加 `skills/<name>`, 重复安装幂等
- `GOMODCACHE` 为空时回退 `go env GOMODCACHE`

## 问题排查

- 是否标注了 `// +skill:<name>`
- 模块是否已在 cache / replace 路径中, 且源目录存在
- 是否在能找到 `go.mod` 的目录下执行

## 参考

- `github.com/xoctopus/genx/pkg/agent`
- `github.com/xoctopus/genx/.agents/skills/genx`
