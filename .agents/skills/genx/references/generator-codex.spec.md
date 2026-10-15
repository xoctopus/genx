# 错误码生成器 (`genx:code`)

生成器标识符: `code` (包 `github.com/xoctopus/genx/devpkg/codex`).

为带 `// +genx:code` 的整型错误码类型生成:

- `Message() string` — 与 `github.com/xoctopus/x/codex` 配合使用
- `StatusCode() int` — **可选**; 与 `httpx.StatusDescriber` (`StatusCode() int`) 对齐, 供 HTTP 层映射

运行时仍用 `codex.New` / `codex.Errorf` 包装枚举值; **不必**改 `x/codex`.

## 接入

在 gen 入口 blank import:

```go
import _ "github.com/xoctopus/genx/devpkg/codex"
```

执行: `genx.Get("code")` 或 `genx.Get()...`.

## 类型与常量命名

类型名 `ErrorCode` → 前缀 `ERROR_CODE` (UpperSnakeCase).

| 常量     | 命名                               |
| -------- | ---------------------------------- |
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
- 可选在 **类型** doc 上覆盖: `// @code domain=CONFWS` (与 `Message` 前缀有关, **不是** HTTP)
- 常量 doc 上 **不要** 写 `domain`; 常量侧 codex 注解只使用 `api` (见下)

参见 `genx/devpkg/testdata/codex/domain_code.go`.

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
| ------------------------- | ------------------------------------------------------------------------  |
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

- 枚举值实现 `StatusCode() int` 即可满足 `httpx.StatusDescriber`
- `codex.New(ERROR__X)` **本身**不是 `StatusDescriber`; HTTP 映射在 httpx `status` 层从 `codex` 取出码再断言 (后续实现)
- 业务复合响应码 (如 `400000001`) 与 `@code api` 无关, 见 `httpx/internal/status.Response`

## 参考

- `github.com/xoctopus/x/codex`
- `github.com/xoctopus/genx/devpkg/codex/doc.go`
