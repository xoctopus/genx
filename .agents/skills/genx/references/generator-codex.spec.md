# 错误码生成器 (`genx:code`)

生成器标识符: `code` (包 `github.com/xoctopus/genx/devpkg/codex`).

为带 `// +genx:code` 的整型错误码类型生成:

| 方法                      | 生成条件                  | 说明                                                    |
|---------------------------|---------------------------|---------------------------------------------------------|
| `Message() string`        | 总是                      | 满足 `codex.ErrorMessage`, 用作错误信息                 |
| `StatusCode() int`        | 任一常量标注 `@code api`  | 供 HTTP 层映射                                          |
| `CodeIdentifier() string` | 总是                      | 错误标识: `<domain>__<NAME>`                            |
| `CodeHint() string`       | 任一常量标注 `@code hint` | 可展示给用户的提示文案                                  |
| `CodeValues() []any`      | 总是                      | 满足 `codex.CanBeCode`, 列出除 `UNDEFINED` 外的全部常量 |

运行时用 `codex.New` / `codex.Wrap` 等包装枚举值. 生成的代码不 import `x/codex` 或 httpx 的接口, 只生成签名匹配的方法.

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

## `Message()` 文案

- 优先取常量 doc 的标题行 (`// ERROR__X human text`)
- 无文案时用 `__` 后的 `<NAME>` 作为 fallback
- 生成格式: `[<domain>:<iota>] <message>`

### 领域码 (`domain`)

- 默认: `UPPER(包名)::UPPER(类型名)`, 例如 `VALIDATION::ERROR`
- 可选在 **类型** doc 上覆盖: `// @code domain=CONFWS` (与 `Message` 前缀和 `CodeIdentifier` 有关, **不是** HTTP)
- 常量 doc 上 **不要** 写 `domain`; 常量侧 codex 注解只使用 `api` 和 `hint` (见下)

参见 `genx/devpkg/testdata/codex/domain_code.go`.

## 错误标识: `CodeIdentifier()`

- 生成规则: `<domain>__<NAME>`, `<NAME>` 为常量名 `__` 之后的部分
- 未定义常量: `<domain>__UNDEFINED`
- 未知值: 返回 `""`

| domain                   | 常量                      | 标识                                  |
|--------------------------|---------------------------|---------------------------------------|
| 默认 `VALIDATION::ERROR` | `ERROR__MISSING_REQUIRED` | `VALIDATION::ERROR__MISSING_REQUIRED` |
| `@code domain=USER`      | `ERROR__NOT_FOUND_BY_ID`  | `USER__NOT_FOUND_BY_ID`               |

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

`0` 表示 "此码不声明 HTTP 语义"; httpx 侧应继续用 Location/Position 或其它规则.

### 校验

- 出现 `@code api=...` 但值非法 → **日志警告**, 跳过该常量的 `StatusCode` case (生成继续)
- 不接受裸 `api=500`; 请用 `StatusInternalServerError` 或 `http=500`

实现通过 `go/importer` 加载当前工具链的 `net/http` 解析 `Status*` 名称.

## 生成文件

- 默认: `{typename}_genx_code.go`
- 嵌入 `AggregationGeneratorMarker` 时: 包内 `zz_genx_code.go`

示例见 `genx/devpkg/testdata/codex/code.go` 与生成结果 `code_genx_code.go`.

## 与 httpx 的关系

- `codex.New(ERROR__X)` **本身**不实现 `StatusCode()` 等方法; httpx 对错误链节点断言 `codex.Valuer`, 取出 code 后再断言 `StatusCode()` / `CodeIdentifier()` / `CodeHint()`
- 生成的 code 类型 **不实现** `error`, 不能直接作为 error 返回, 需经 `codex.New` 等包装
- 业务复合响应码 (如 `400000001`) 与 `@code api` 无关, 见 `httpx/internal/status.Response`

## 参考

- `github.com/xoctopus/x/codex`
- `github.com/xoctopus/genx/devpkg/codex/doc.go`
