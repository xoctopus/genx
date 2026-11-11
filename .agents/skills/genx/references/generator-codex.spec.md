# 错误码生成器 (`genx:code`)

生成器标识符: `code` (包 `github.com/xoctopus/genx/devpkg/codex`).

为带 `// +genx:code` 的整型错误码类型生成:

| 方法                      | 生成条件                  | 说明                                                       |
|---------------------------|---------------------------|------------------------------------------------------------|
| `Message() string`        | 总是                      | 满足 `codex.HasMessageCode`, 用作错误信息                  |
| `StatusCode() int`        | 任一常量标注 `@code api`  | 错误码的 HTTP 状态码语义                                   |
| `CodeIdentifier() string` | 总是                      | 满足 `codex.HasIdentifierCode`, 错误标识 `<DOMAIN>::<KEY>` |
| `CodeHint() string`       | 任一常量标注 `@code hint` | 满足 `codex.HasHintCode`, 可展示给用户的提示文案           |
| `CodeValues() []any`      | 总是                      | 满足 `codex.CanBeCode`, 列出除 `UNDEFINED` 外的全部常量    |

运行时用 `codex.New` / `codex.Wrap` 等包装枚举值. 生成的代码不 import 任何接口定义, 只生成签名匹配的方法.

## 接入

在 gen 入口 blank import:

```go
import _ "github.com/xoctopus/genx/devpkg/codex"
```

执行: `genx.Get("code")` 或 `genx.Get()...`.

## 类型与常量命名

类型名 `ErrorCode` → 前缀 `ERROR_CODE` (UpperSnakeCase).

| 常量     | 命名                               |
|----------|------------------------------------|
| 未定义   | `<PREFIX>_UNDEFINED`               |
| 具体错误 | `<PREFIX>__<NAME>` (双下划线 `__`) |

示例:

```go
// +genx:code
type Error int8

const (
ERROR_UNDEFINED Error = iota
ERROR__NOT_FOUND // user not found
)
```

### 数值偏移 (offset)

错误码数值只在本类型内唯一, 不同包的错误码数值可能相同 (都从 `1` 开始). 需要数值在项目内也唯一时, 给常量加上偏移:

```go
// DomainOffset 由项目集中分配, 保证各包不重叠
const DomainOffset = Error(100000)

// +genx:code
type Error int

const (
	ERROR_UNDEFINED Error = iota + DomainOffset // 100000
	ERROR__NOT_FOUND                            // 100001
)
```

| 约束             | 说明                                                     |
|------------------|----------------------------------------------------------|
| 偏移值由项目管理 | 生成器不分配, 也不检查跨包重叠                           |
| 底层类型要够宽   | 如 `int8` 放不下 `100000`, 需改为 `int` / `int32`        |
| 偏移常量的命名   | 不能符合 `<PREFIX>__<NAME>`, 否则会被当作错误码生成      |
| 生成结果         | `Message()` 数值带偏移, 如 `[<domain>:100001]`; 标识不变 |

偏移只解决数值重叠; 错误标识 (`CodeIdentifier()`) 由 domain 区分, 不依赖数值.

参见 `genx/devpkg/testdata/codex/code.go`.

## `Message()` 文案

- 优先取常量 doc 的标题行 (`// ERROR__X human text`)
- 无文案时用 `__` 后的 `<NAME>` 作为 fallback
- 生成格式: `[<domain>:<iota>] <message>`

### 领域码 (`domain`)

| 场景                      | `Message()` 前缀             | `CodeIdentifier()` 的 domain |
|---------------------------|------------------------------|------------------------------|
| 默认                      | `UPPER(包名)::UPPER(类型名)` | `UPPER(包名)`                |
| 类型 doc `@code domain=X` | `X`                          | `X`                          |

- 例如包 `validation` 中的 `Error`: 前缀 `VALIDATION::ERROR`, 标识 domain `VALIDATION`
- 覆盖只写在 **类型** doc 上: `// @code domain=CONFWS` (与 HTTP 无关)
- 常量 doc 上 **不要** 写 `domain`; 常量侧 codex 注解只使用 `api` 和 `hint` (见下)
- 推荐一个包只定义一个错误码类型 (包即领域). 多个类型可以生成, 但不推荐: 它们的标识共用同一个 domain

参见 `genx/devpkg/testdata/codex/domain_code.go`.

## 错误标识: `CodeIdentifier()`

- 生成规则: `<DOMAIN>::<KEY>`, `<KEY>` 为常量名 `<PREFIX>__` 之后的部分, `<DOMAIN>` 见上表
- 未定义常量: `<DOMAIN>::UNDEFINED`
- 未知值: 返回 `""`

| 包 / 注解                      | 常量                      | 标识                           |
|--------------------------------|---------------------------|--------------------------------|
| 包 `validation`, 无 `domain`   | `ERROR__MISSING_REQUIRED` | `VALIDATION::MISSING_REQUIRED` |
| 包 `user`, `@code domain=USER` | `ERROR__NOT_FOUND_BY_ID`  | `USER::NOT_FOUND_BY_ID`        |
| 包 `user`, `@code domain=USER` | `ERROR_UNDEFINED`         | `USER::UNDEFINED`              |

## 提示文案: `@code hint`

在 **具体错误常量** 的 doc 中标注, 用于生成 `CodeHint()`. 生成条件与 `StatusCode()` 一致.

| 写法              | 含义                                                        |
|-------------------|-------------------------------------------------------------|
| `@code hint=...`  | 使用标注内容作为提示文案                                    |
| `@code hint`      | 使用 doc 标题行作为提示文案; 标题为空时 **日志警告** 并跳过 |
| 未写 `hint`       | 该常量不提示; `CodeHint()` 走 `default` → **`""`**          |
| 整型无任何 `hint` | **不生成** `CodeHint()` 方法                                |

文案不带 `[<domain>:<value>]` 前缀.

```go
// +genx:code
type Error int8

const (
	ERROR_UNDEFINED Error = iota
	// ERROR__NOT_FOUND user not found
	// @code hint
	ERROR__NOT_FOUND
	// ERROR__INVALID_NAME invalid name
	// @code hint=name must be 8-16 letters or digits
	ERROR__INVALID_NAME
	ERROR__INTERNAL // internal error
)
```

生成 `ERROR__NOT_FOUND.CodeHint() == "user not found"`, `ERROR__INVALID_NAME.CodeHint() == "name must be 8-16 letters or digits"`, `ERROR__INTERNAL.CodeHint() == ""`.

## HTTP 状态: `@code api=...` (仅常量)

在 **具体错误常量** 的 doc 中标注, 用于生成 `StatusCode()`:

```go
// ERROR__NOT_FOUND user not found
// @code api=StatusNotFound
ERROR__NOT_FOUND

// ERROR__CLIENT_CLOSED client closed (non-IANA)
// @code api=http=499
ERROR__CLIENT_CLOSED
```

### 取值规则

| 写法                      | 含义                                                                      |
|---------------------------|---------------------------------------------------------------------------|
| `api=StatusBadRequest` 等 | `net/http` 中以 `Status` 开头的常量名; 生成 `return 400` 等数值           |
| `api=http=444`            | 三位 HTTP 码 `100`-`999`; 用于 IANA/`net/http` 未收录的码 (如 444)        |
| 未写 `api`                | 该常量不参与 HTTP 映射; `StatusCode()` 的 `switch` 走 `default` → **`0`** |
| 整型无任何 `api`          | **不生成** `StatusCode()` 方法                                            |

`0` 表示 "此码不声明 HTTP 语义", 如何处理由使用方决定.

### 校验

- 出现 `@code api=...` 但值非法 → **日志警告**, 跳过该常量的 `StatusCode` case (生成继续)
- 不接受裸 `api=500`; 请用 `StatusInternalServerError` 或 `http=500`

实现通过 `go/importer` 加载当前工具链的 `net/http` 解析 `Status*` 名称.

## 生成文件

- 默认: `{typename}_genx_code.go`
- 嵌入 `AggregationGeneratorMarker` 时: 包内 `zz_genx_code.go`

示例见 `genx/devpkg/testdata/codex/code.go` 与生成结果 `code_genx_code.go`.

## 运行时使用

- 生成的 code 类型 **不实现** `error`, 需经 `codex.New` / `codex.Errorf` / `codex.Wrap` / `codex.Wrapf` 包装后作为 error 返回
- 不知道 code 类型的使用方, 对错误链节点断言 `codex.Error`, 用 `Error.As` 把节点视为 code, 再断言上述方法 (如 `StatusCode()`, `codex.HasIdentifierCode`)

## 参考

- `github.com/xoctopus/x/codex`
- `github.com/xoctopus/genx/devpkg/codex/doc.go`
