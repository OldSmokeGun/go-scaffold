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

### Adapter 只做协议转换，不编排业务逻辑

Adapter 的职责边界可以概括为一句话：**它只负责"外部协议 ↔ Controller 入口"的双向转换**。

* **入方向**：把外部协议表示（HTTP 请求体、gRPC 消息、CLI 参数、cron 配置等）解析、校验格式、转换为 Controller 的业务入参结构，然后调用 Controller。
* **出方向**：把 Controller 返回的业务结果（或业务错误）转换为外部协议表示（HTTP 响应、gRPC 响应、退出码等），然后返回给调用方。

Adapter 在整个调用链中只扮演"翻译"角色：

```text
            入方向                         出方向
外部协议 ──解析/DTO 转换──► Controller 入参   Controller 结果 ──DTO 转换/渲染──► 外部协议
```

Adapter **不编排业务逻辑**，具体来说，Adapter 不得：

* 决定调用哪个 Usecase、以什么顺序调用多个业务操作。
* 在 handler 内组合多个 Controller 调用来"拼装"业务结果。
* 对 Controller 的返回结果做业务语义上的加工（合并、拆分、条件取舍）。
* 根据业务含义做决策（例如"库存不足就不创建订单"这类判断）。

禁止的示例：

```go
func (h *Handler) CreateOrder(c echo.Context) error {
	var req CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 禁止：handler 内编排多个业务操作
	if err := h.controller.CheckInventory(c.Request().Context(), req.ProductID, req.Quantity); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := h.controller.CheckBalance(c.Request().Context(), req.UserID, ...); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	result, err := h.controller.CreateOrder(c.Request().Context(), input)
	...
}
```

要求：编排发生在 Controller 内，Adapter 只调用一次 Controller 工作流入口：

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

"步骤如何组织"属于业务工作流，必须由 Controller 编排；Adapter 永远只关心"这个协议长什么样"。

### Controller 不感知外部协议

与上一条对称，**Controller 必须完全不知道外部协议的存在**：

* Controller 的入参/出参必须是业务语义结构（如 `CreateOrderInput`、`*domain.Order`），不得是 `*http.Request`、`echo.Context`、gRPC 消息等协议类型。
* Controller 不设置 HTTP 状态码、不构造 HTTP/gRPC 响应体、不感知路由或序列化格式。
* 同一个 Controller 工作流应当可以同时被 HTTP、gRPC、cron、CLI 等任意 Adapter 复用，而无需修改一行代码。

```text
HTTP Adapter     ─► Controller.CreateOrder ─► HTTP 响应
gRPC Adapter     ─► Controller.CreateOrder ─► gRPC 响应
Cron Job         ─► Controller.CreateOrder ─► （忽略结果）
```

如果 Controller 的签名中出现了协议类型，说明协议概念已经泄漏进业务层，必须回退到 Adapter 转换。

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
