# 错误码类型定义指南

## 错误码类型接口定义

参考 `github.com/xoctopus/x/codex`

## 生成器接入方式

按照 `github.com/xoctopus/genx/devpkg/codex` 约定接入

## 错误码的领域码

默认使用 `{包名}.{类型名}` 作为默认领域码

使用 `@code domain=...` 自定义领域码

参见 `github.com/xoctopus/genx/devpkg/testdata/codex/domain_code.go`