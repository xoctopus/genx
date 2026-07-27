---
name: genx-guideline
description: 封装 genx 的自定义生成器扩展方式与项目接入约定; 当任务涉及生成器扩展, 注册, 组装执行入口或排查 +genx 生成行为时使用.
---

# genx-guideline

按 `github.com/xoctopus/genx` 约定接入代码生成或扩展自定义生成器。

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
	entry = ""
)

func main() {
	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry},
	})

	must.NoError(ctx.Execute(context.Background(), genx.Get()...))
}
```

**关键约定**：

- 定义约定: 定义生成器包, 并通过 `init()` 方式注册生成器
- 生成目标: `Entrypoint` 指定需要生成的包目录或包路径
- 输出目标: 如果实现 `genx.AggregationGeneratorMarker` 每个包的生成文件, 统一输出到 `zz_genx_{identifier}.go`; 否则按每个类型输出到 `{typename}_genx_{identifier}.go`

## 更多信息

参见 [genx.spec.md](references/genx.spec.md) 

