# 无符号整数文本序列化生成器指南

## 功能描述

为无符号整数类型 (`uint`, `uint8`, `uint16`, `uint32`, `uint64`) 自动生成 `encoding.TextMarshaler` 和 `encoding.TextUnmarshaler` 接口的实现.

这允许基于无符号整数的自定义类型 (例如 `type UserID uint64`) 在 JSON, XML 或纯文本等格式中, 以十进制字符串的形式进行序列化和反序列化.

## 生成器接入方式

按照 `github.com/xoctopus/genx/devpkg/uintstr` 约定接入.

在需要生成的类型上方添加 `// +genx:uintstr` 指令:

```go
// +genx:uintstr
type UserID uint64
```
