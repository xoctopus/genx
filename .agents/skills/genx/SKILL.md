---
name: genx-guideline
description: >-
  封装 genx 的自定义生成器扩展方式与项目接入约定, 以及 skills 安装.
  当任务涉及生成器扩展, 注册, 组装执行入口, 或排查 +genx 生成行为时使用.
triggers:
  - "添加/扩展 genx 生成器"
  - "接入 genx 代码生成"
  - "排查 genx 未生成"
  - "安装 skills"
---

# genx-guideline

按 `github.com/xoctopus/genx` 约定接入代码生成或扩展自定义生成器.

## 接入项目

```go
import (
	_ "github.com/xoctopus/genx/devpkg/contextx"
	_ "github.com/xoctopus/genx/devpkg/docx"
	_ "github.com/xoctopus/genx/devpkg/enumx"
	_ "github.com/xoctopus/genx/devpkg/codex"
	_ "github.com/xoctopus/genx/devpkg/uintstr"

	"github.com/xoctopus/genx/pkg/genx"
)

var (
	// eg: ./... github.com/to/pkg/path ./pkg/foo /abs/path
	entry = "local path or package path"
)

func main() {
	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry},
	})

	must.NoError(ctx.Execute(context.Background(), genx.Get()...))
}
```

如果执行全部生成器

`ctx.Execute(context.Background(), genx.Get()...)`

如果需要指定某些生成器

`ctx.Execute(context.Background(), genx.Get("identifier1", "identifier2")...)`

- `Get(id)` 决定哪些生成器被执行
- 代码中 `// +genx:<id>` 注释决定哪些类型或包被生成
- 如果 `<id> == Generator.Identifier()` 会命中生成器生成代码

**关键约定**:

- 定义约定: 定义生成器包, 并通过 `init()` 方式注册生成器
- 生成目标: `Entrypoint` 指定需要生成的包目录或包路径
- 输出目标: 如果嵌入 `genx.AggregationGeneratorMarker` 每个包的生成文件, 统一输出到 `zz_genx_{identifier}.go`; 否则按每个类型输出到 `{typename}_genx_{identifier}.go`

## 内置生成器

- `code`: 错误码. 参见 [generator-codex.spec.md](references/generator-codex.spec.md)
- `context`: 上下文注入. 参见 [generator-contextx.spec.md](references/generator-contextx.spec.md)
- `doc`: 运行时文档. 参见 [generator-docx.spec.md](references/generator-docx.spec.md)
- `enum`: 枚举. 参见 [generator-enumx.spec.md](references/generator-enumx.spec.md)
- `uintstr`: 无符号整型序列号反序列化. 参见 [generator-uintstr.spec.md](references/generator-uintstr.spec.md)

## 实现参考

参见

- `github.com/xoctopus/genx/devpkg/codex`
- `github.com/xoctopus/genx/devpkg/contextx`
- `github.com/xoctopus/genx/devpkg/docx`
- `github.com/xoctopus/genx/devpkg/enumx`
- `github.com/xoctopus/genx/devpkg/uintstr`

## 代码生成底座

- 如果是与类型无关纯文本代码生成可以考虑使用无参数模版 `snippet.Template`
- 如果是与类型相关的代码生成, 不强制使用 `snippet` 包, 但是需要基于语法分析.
- 文档的 `描述`, `标注`, `指令` 需要遵循 `github.com/xoctopus/genx/pkg/docx` 规则

## 问题排查

- 是否导入生成器. 比如: `import _ pkg/to/your/generator`
- 入口是否配置 `Entrypoint`. `Entrypoint` 指定需要生成的包目录或包路径
- 生成器是否注册. 是否在 `init()` 调用了 `genx.Register`
- 生成目标是否已经定义在了包注释或类型定义注释中. (如果定义在package doc 那么该生成器会应用到包内所有类型定义. 比如 `+genx:doc`)

## 更多信息

参见 [genx.spec.md](references/genx.spec.md)

如果需要在项目中集成 skills 安装

参见 [skills-installation.spec.md](references/skills-installation.spec.md)
