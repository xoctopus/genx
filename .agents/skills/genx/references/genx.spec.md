# 自定义代码生成器规范

本文描述如何基于 `genx` 自定义代码生成器.

## 适用范围

适用一下任务:

- 自定义代码生成器
- 修改已有生成器
- 排查生成器执行异常

## 主入口

### 运行时接口

- `github.com/xoctopus/genx/pkg/genx`
  - `Generator`: 生成器接口
  - `Versioned`: 生成器版本信息
  - `AggregationGeneratorMarker`: 聚合式生成器接口. 如果实现改接口则会将生成文件聚合到同一个生成文件, 否则会按照类型分文件输出.
  - `GeneratorNewer`: 按照 `Context` 创建 `Generator` 实例

### 参考实现

- `github.com/xocotpus/genx/devpkg/enumx`
- `github.com/xocotpus/genx/devpkg/codex`
- `github.com/xocotpus/genx/devpkg/contextx`
- `github.com/xocotpus/genx/devpkg/docx`
- `github.com/xocotpus/genx/devpkg/uintstr`

## 自定义方式

### 注册约定

通过 `init()` 方式注册生成器:

1. 在独立包中实现生成器
2. 在 `init()` 中调用 `genx.Register(...)`

### 生成器接口

- `Generator.Identifier` 生成器唯一名称标识
- `Generator.Generate` 实现 `types.Type` 类型的代码生成

### 文档解析规范

文档解析规范定义实现在 `github.com/xoctopus/genx/pkg/docx`

- 注释首行为文档标题或概要描述
- 注释如果以 `+` 起始, 则是生成器生成指令标识. 如: `+genx:enum`
- 注释如果以 `@` 起始, 则是生成器附加描述. 后面紧跟生成器标识, 表示针对某个生成器的附加描述. 如 `@enum storage=text`
- 其余非特殊注释为文档正文, 做更详细的描述

举例:

```golang
package demo

// Status 启用状态
// 开启关闭二元状态描述
// +genx:enum
// @enum storage=int
type Status int
```

以上面代码片段为例:

类型 `Status`

- 文档标题: `启用状态`
- 文档正文: `开启关闭二元状态描述`
- enum生成器标识: `+genx:enum`
- enum生成器标注: `@enum storage=int`

生成器实现时, 根据类型或值的文档做为代码生成依据.

### 输出约定

- 生成器实现 `AggregationGeneratorMarker` 则某个包内到代码生成统一输出到同一文件. 格式: `zz_genx_{identifier}.go`
- 否则按照类型分别输出到对应文件. 格式: `{typename}_genx_{identifier}.go`

## 如何接入

在包级别注释, 类型注释中添加生成指令. 如: `+genx:doc`

## 问题排查

- 入口是否配置 `Entrypoint`. `Entrypoint` 指定需要生成的包目录或包路径
- 生成器是否注册. 是否在 `init()` 调用了 `genx.Register`
- 生成目标是否已经定义在了包注释或类型定义注释中.
