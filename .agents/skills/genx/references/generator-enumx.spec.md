# 枚举类型定义指南

## 枚举类型接口定义

参考 `github.com/xoctopus/x/enumx`

## 生成器接入方式

按照 `github.com/xoctopus/genx/devpkg/enumx` 约定接入

```shell
go doc github.com/xoctopus/genx/devpkg/enumx
```

## 选项指南

### 存储类型

如果 DB 按文本存储, 需要标注 `storage`

参见 `github.com/xoctopus/genx/devpkg/testdata/enumx/gender.go`

`@enum storage=...` 会影响数据库值的序列化和反序列化.

storage 默认: `int`
storage 可选: `text`, `string`, `varchar`, `enum`

### 枚举类型的映射关系

如果需要枚举类型和另外的枚举类型的映射关系, 则需要在类型定义注释中标注映射关系.

格式为 `@enum map.{映射方法名}={映射类型名}` 这里的映射类型必须是同一个包内的类型. 如果是其他包需要自行提前定义别名.

并在需要的枚举值注释中标注具体的映射关系.

格式为 `@enum map.{映射方法名}={映射值变量名}`

参见 `github.com/xoctopus/genx/devpkg/testdata/enumx/product_type.go`

### 枚举类型的扩展方法

如果需要枚举类型更多的映射方法, 则需要枚举值注释中标注具体的映射值. (映射方法返回值仅支持 `string`)

格式为 `@enum ext.{映射方法名}={字符串}`

注意 `string`, `text`, `value` 是 `@enum ext.` 保留字, 不会生成.

参见 `github.com/xoctopus/genx/devpkg/testdata/enumx/gender.go`

注意

- `@enum map.` 映射其他枚举类型
- `@enum ext.` 映射文本(string)
