# 错误处理

位置：`internal/errors`

---

## 职责

`errors` 包定义标准化的应用错误和业务错误。

它可以包含：

* 业务错误码
* 业务错误消息
* 与业务错误关联的 HTTP 状态码
* 标准应用错误
* 错误构造函数
* 预定义的业务错误

示例：

```go
var ErrInsufficientBalance = BizError(
	40001,
	400,
	"余额不足",
)
```

---

## 错误归属

业务错误必须集中定义。

不得在多个 Usecase、Controller 或 Adapter 中重复定义同一个业务错误。

禁止：

```go
return errors.New("余额不足")
```

当已存在对应的集中定义的业务错误时。

要求：

```go
return errors.ErrInsufficientBalance
```

---

## errors 层限制

`errors` 包不得：

* 访问 repository 或数据库。
* 调用 Usecase 或 Controller。
* 处理 HTTP 请求。
* 执行业务工作流。
* 依赖传输层实现。

Adapter 可以将业务错误转换为协议特定的表示：

```text
internal/errors.ErrInsufficientBalance
              │
              ▼
HTTP Adapter
              │
              ▼
HTTP 400 + 响应体
```

HTTP 状态码映射不得在 Usecase 内实现。

---

## 错误传播

只要需要添加额外上下文，错误就必须保留其原始原因。

使用：

```go
return fmt.Errorf("query goods info: %w", err)
```

不要使用：

```go
return fmt.Errorf("query goods info: %v", err)
```

当原始错误需要保持可检查时。

不得静默丢弃错误。

* 当不需要额外上下文时，`internal/errors` 中定义的业务错误应直接返回。
* 基础设施错误在返回给调用方之前，通常应附上有意义的上下文进行包装。
