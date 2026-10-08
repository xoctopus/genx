# 业务仓接入手册

本文说明 **业务仓库** (可执行服务, 多 cmd, 依赖 confx/sqlx 等) 如何接入:

1. **代码生成器** (`internal/cmd/gen`, genx + 各库 `devpkg`)
2. **Skill 安装器** (`internal/cmd/skill-install`, genx `pkg/agent`)

工程脚手架 (Makefile, Docker, CI) 由 **devx** / `devgen` 负责, 不在本文范围.

## 1. 总览

```text
go.mod
  require  + // +skill:<name>  →  声明依赖与要安装的 Agent skill
  tool     →  本仓 internal/cmd/gen, internal/cmd/skill-install

internal/cmd/gen/main.go           →  blank import devpkg, 执行 genx
internal/cmd/skill-install/main.go →  agent.Installer 安装 .agents/skills

源码 // +genx:<id> 等指令          →  生成目标
go run ./internal/cmd/gen          →  产出 *_genx_*.go 等
go run ./internal/cmd/skill-install →  symlink 各模块 skill
```

| 入口                         | 作用                   | 依赖                       |
|------------------------------|------------------------|----------------------------|
| `internal/cmd/gen`           | 跑已注册的 genx 生成器 | `genx` + 各能力库 `devpkg` |
| `internal/cmd/skill-install` | 按 `go.mod` 安装 skill | 仅 `genx/pkg/agent`        |

二者独立: 装 skill 不触发代码生成; 跑 gen 不安装 skill.

## 2. go.mod

### 2.1 require 与 +skill

在
**直接** `require` 上一行标注 skill. 名称必须与提供方目录 `.agents/skills/<name>/` 一致 (不是 SKILL frontmatter 里的 `name` 字段).

```go
module example.com/myapp

go 1.27.0

tool (
	example.com/myapp/internal/cmd/gen
	example.com/myapp/internal/cmd/skill-install
)

require (
	// +skill:genx
	github.com/xoctopus/genx v0.3.9
	// +skill:sqlx
	github.com/xoctopus/sqlx v0.4.6
	// +skill:concx
	github.com/xoctopus/confx v0.6.2
	// +skill:testx
	github.com/xoctopus/x v0.5.9
)
```

- 仅 **direct require** 参与 skill 安装; `indirect` 忽略.
- 本地联调可用 `replace` 指向 monorepo 路径 (见 [skills-installation.spec.md](skills-installation.spec.md)).
- `tool` 块 **不会** 自动安装 skill; 需要 skill 的模块应出现在 `require` 且带 `+skill`.

### 2.2 首次整理

```bash
go mod tidy
```

## 3. Skill 安装器

### 3.1 创建 `internal/cmd/skill-install/main.go`

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/xoctopus/genx/pkg/agent"
)

func main() {
	if err := (&agent.Installer{}).Install(context.Background()); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

### 3.2 用法

在 **模块根目录** (含 `go.mod`) 执行:

```bash
go run ./internal/cmd/skill-install
# 或
go tool skill-install   # 若已通过 go.mod tool 安装到 PATH
```

结果:

- 创建/更新 `.agents/skills/<name>` → 指向依赖模块内 `.agents/skills/<name>` 的 symlink
- 在 `.agents/.gitignore` 追加 `skills/<name>` (幂等)

### 3.3 排查

- `+skill` 名称与远端目录名不一致 → 安装失败或目录不存在
- 未在模块根执行 → 找不到 `go.mod`
- 依赖模块版本在 cache 中无 `.agents/skills` → 升级版本或 `replace` 到源码树

协议细节: [skills-installation.spec.md](skills-installation.spec.md).

## 4. 代码生成器

### 4.1 创建 `internal/cmd/gen/main.go`

**推荐**: 一次引入 genx 内置生成器:

```go
package main

import (
	"context"
	"os"

	_ "github.com/xoctopus/genx/devpkg"
	"github.com/xoctopus/genx/pkg/genx"
)

func main() {
	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{"./..."},
	})

	if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
		panic(err)
	}
}
```

**按需** 增加其它库的 devpkg (未 import 则对应生成器不会注册):

```go
import (
_ "github.com/xoctopus/genx/devpkg"
_ "github.com/xoctopus/sqlx/devpkg" // sqlx model 等

"github.com/xoctopus/genx/pkg/genx"
)
```

**额外扫描路径** (如独立 models 目录, 不在 `./...` 下时):

```go
import (
"context"
"os"
"path/filepath"

_ "github.com/xoctopus/genx/devpkg"
"github.com/xoctopus/genx/pkg/genx"
_ "github.com/xoctopus/sqlx/devpkg"
)

func main() {
cwd, _ := os.Getwd()
ctx := genx.NewContext(&genx.Args{
Entrypoint: []string{
"./...",
filepath.Join(cwd, "pkg", "models"),
},
})
if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
panic(err)
}
}
```

参考: `github.com/xoctopus/confx/internal/cmd/gen/main.go` (含 `pkg/models` 等额外 Entrypoint).

### 4.2 在源码中声明生成目标

- 包或类型注释: `// +genx:enum`, `// +genx:doc`, `// +genx:code` 等
- sqlx 等库自有指令见对应 skill (如 `// +genx:model` 以 sqlx 文档为准)

### 4.3 用法

```bash
go run ./internal/cmd/gen
# 或
go tool gen
```

只跑部分生成器:

```go
ctx.Execute(context.Background(), genx.Get("enum", "doc")...)
```

- `Get(id...)`: 本次执行哪些 **已注册** 生成器
- `+genx:<id>`: 源码里哪些包/类型需要生成

### 4.4 与「逐个 import devpkg」的关系

SKILL 中的最小示例逐个 `_ "…/devpkg/enumx"` 等价于 `_ "github.com/xoctopus/genx/devpkg"`.
业务仓优先用 **devpkg 聚合 import**, 再叠加 sqlx 等 **外部 devpkg**.

自定义生成器: `import _ "your/module/yourgen"` 并在该包 `init()` 里 `genx.Register`.

## 5. 推荐接入顺序 (新业务仓)

1. `go.mod`: 声明 `require` + `// +skill:*`, 添加 `tool` 两个 internal cmd
2. 编写 `internal/cmd/skill-install/main.go`, 执行 `go mod tidy` 后 `go run ./internal/cmd/skill-install`
3. 编写 `internal/cmd/gen/main.go` (genx devpkg + 需要的 sqlx 等)
4. 在业务代码添加 `+genx` / sqlx 指令
5. `go run ./internal/cmd/gen`, 将生成文件纳入版本管理
6. (可选) `devgen init` / `devgen all` 生成 Makefile 与 CI — 见 devx 文档

## 6. 参考仓库文件

| 仓库  | gen                        | skill-install                        |
|-------|----------------------------|--------------------------------------|
| confx | `internal/cmd/gen/main.go` | `internal/cmd/skill-install/main.go` |
| sqlx  | `internal/cmd/gen/main.go` | `internal/cmd/skill-install/main.go` |
| devx  | `internal/cmd/gen/main.go` | `internal/cmd/skill-install/main.go` |

## 7. 相关文档

- 自定义生成器实现: [genx.spec.md](genx.spec.md)
- Skill 安装协议: [skills-installation.spec.md](skills-installation.spec.md)
- 总览与排查: 上级 [SKILL.md](../SKILL.md)
