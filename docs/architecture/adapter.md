# Adapter 层

位置：`internal/app/adapter`

---

## 职责

Adapter 层是**外部协议转换层**，也是应用的外部入口。

Adapter 可能包括：

* HTTP
* gRPC
* cron
* CLI
* 脚本
* 其他外部调用机制

Adapter 将外部表示转换为内部应用表示，并将内部结果转换为外部表示。

示例流程：

```text
HTTP 请求
    ↓
HTTP handler
    ↓
Controller 请求
```

```text
Controller 结果
    ↓
HTTP 响应
```

Adapter 层不得包含业务逻辑。

---

## Adapter 可以做的事

* 解析 HTTP 请求参数 / 请求体。
* 解析 gRPC 请求消息。
* 解析 CLI 参数。
* 解析 cron/任务配置。
* 执行协议特定的校验。
* 将外部 DTO 转换为 Controller 的输入结构。
* 将 Controller 的输出转换为 HTTP/gRPC/CLI 响应结构。
* 设置 HTTP 状态码和响应头。
* 将业务错误转换为协议特定的错误表示。
* 处理协议特定的认证/授权中间件。
* 处理序列化 / 反序列化。
* 处理协议特定的日志与追踪。

正确的示例：

```go
func (h *Handler) CreateOrder(c echo.Context) error {
	var req CreateOrderRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	input := controller.CreateOrderInput{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	result, err := h.controller.CreateOrder(c.Request().Context(), input)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, toResponse(result))
}
```

---

## Adapter 不得做的事

* 直接查询 MySQL / Redis。
* 作为业务处理的一部分直接发布消息。
* 调用 Repository 方法。
* 实现业务规则或工作流。
* 计算业务价格 / 检查库存 / 检查余额。
* 创建业务订单或以业务决策的方式修改业务实体。
* 编排多个 Usecase。
* 决定业务结果。

强制调用路径：

```text
Adapter → Controller          // 允许
Adapter → Usecase             // 禁止
Adapter → Repository          // 禁止
Adapter → Database / Redis    // 禁止
```

即使绕过 Controller 看起来更简单，Adapter 也不得这样做。

---

## Adapter 子目录

```text
adapter
├─cron
│  ├─job
│  └─scheduler
├─scripts
└─server
   ├─grpc
   │  ├─handler
   │  └─router
   └─http
      ├─handler
      ├─middleware
      └─router
```

### `router`

* 路由 / 端点注册。
* 协议级配置。
* 不得包含业务逻辑。

### `handler`

* 接收外部请求。
* 解析输入。
* 调用 Controller。
* 将 Controller 的输出转换为协议响应。
* 不得直接调用 Usecase 或 Repository。

### `middleware`

协议级的横切关注点：

* 认证
* 请求追踪
* 日志
* CORS
* 限流
* 请求元数据

Middleware 不得实现业务工作流。

### `cron/job`

定时触发的外部入口。必须调用 Controller，而不是 Usecase 或 Repository。

### `scheduler`

调度并触发任务。不得包含业务逻辑。

### `scripts`

外部入口。必须遵守与 Adapter 相同的规则，且不得绕过 Controller。

---

## DTO 规则

传输层特定的 DTO 必须保留在 Adapter 内部，例如：

```text
adapter/http/handler
    CreateOrderRequest
    CreateOrderResponse
```

这些类型不得用作业务层实体。

不得将 `*http.Request`、`echo.Context`、gRPC 请求/响应类型传入 Controller 或 Usecase。

应优先使用内部请求结构：

```go
type CreateOrderInput struct {
	UserID    int64
	ProductID int64
	Quantity  int
}
```

由 Adapter 将外部数据转换为该结构。
