# 运行时文档生成器指南

## 运行时文档接口定义

参考
- `github.com/xoctopus/x/docx`: 运行时文档接口定义

## 生成器接入方式

按照 `github.com/xoctopus/genx/devpkg/docx` 约定接入

在需要生成的 package doc 添加 `// +genx:doc` 指令:

```go
// Package xxx ...
// +genx:doc
package xxx
```
