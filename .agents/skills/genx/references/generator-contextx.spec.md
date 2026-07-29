# 上下文(context)使用指南

## 上下文定义

参考 `github.com/xoctopus/x/contextx`

## 生成器接入方式

按照 `github.com/xoctopus/genx/devpkg/contextx` 约定接入

在需要生成的类型上方添加 `// +genx:context` 指令:

```go
// +genx:context
type T struct {...}
```
