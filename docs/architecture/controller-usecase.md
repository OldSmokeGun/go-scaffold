# Controller 与 Usecase

---

## Controller

位置：`internal/app/controller`

### 职责

Controller 是**应用业务的核心入口**。

它负责：

1. 校验业务输入参数。
2. 组织业务操作。
3. 编排多个 Usecase。
4. 定义 Usecase 的执行顺序。
5. 协调不同的业务模块。
6. 控制应用级工作流。
7. 返回最终业务结果。

Controller 不得实现各个模块的内部业务逻辑。

* Controller 定义**需要发生哪些业务操作以及以什么顺序发生**。
* Usecase 定义**每个具体业务操作如何实现**。

### Controller 不得访问 Repository

Controller 不得直接调用 Repository、数据库、Redis、消息队列或外部服务。

Controller 必须通过 Usecase 与业务功能交互。

> 例外：满足「极简单单资源 CRUD」条件时，Controller 可直接调用 Repository 抽象，见下方小节 [例外：极简单单资源 CRUD 可直接调用 Repository](#例外极简单单资源-crud-可直接调用-repository)。

禁止：

```go
func (c *Controller) CreateOrder(...) error {
	goods, err := c.goodsRepository.QueryGoodsInfo(...)
	...
}
```

要求：

```go
func (c *Controller) CreateOrder(...) error {
	if err := validateInput(...); err != nil {
		return err
	}

	if err := c.goodsUsecase.CheckInventory(...); err != nil {
		return err
	}

	if err := c.userUsecase.CheckBalance(...); err != nil {
		return err
	}

	order, err := c.orderUsecase.Create(...)
	if err != nil {
		return err
	}

	return order
}
```

### 例外：极简单单资源 CRUD 可直接调用 Repository

当操作**同时满足**以下全部条件时，Controller 可以直接调用 Repository，以避免引入只是透传参数的空壳 Usecase：

1. 仅涉及单一资源/单一模块。
2. 不涉及跨模块编排或跨模块数据依赖。
3. 除基本参数校验外没有实质性业务规则（无库存扣减、余额计算、定价、权限判定等）。
4. 操作能由单个（或一组简单顺序的）Repository 方法直接表达。

推荐写法（以极简单的字典管理为例）：

```go
func (c *DictController) CreateDict(
	ctx context.Context,
	input CreateDictInput,
) (*domain.Dict, error) {
	if input.Name == "" {
		return nil, errors.ErrInvalidParam
	}

	return c.dictRepository.Create(ctx, repository.CreateDictParams{
		Name:  input.Name,
		Value: input.Value,
	})
}
```

不推荐为这类操作引入纯透传的样板代码：

```go
// Usecase 只是转发，没有任何业务价值 —— 不推荐
func (u *DictUsecase) Create(ctx context.Context, input CreateDictInput) (*domain.Dict, error) {
	return u.dictRepository.Create(ctx, input)
}

// Controller 也只是转发 —— 不推荐
func (c *DictController) CreateDict(ctx context.Context, input CreateDictInput) (*domain.Dict, error) {
	return c.dictUsecase.Create(ctx, input)
}
```

注意事项：

* 该例外仅限 controller → repository 这一步；adapter 仍然必须调用 Controller，不得绕过 Controller。
* 一旦操作将来需要增加业务规则、跨模块步骤或多步编排，必须立即将逻辑下沉到 Usecase，恢复标准的 controller → usecase → repository 路径。
* 该例外属于有意的设计取舍，不得作为将复杂逻辑堆进 Controller 的理由。

### 工作流示例

```text
CreateOrder
    │
    ├── GoodsUsecase.CheckInventory
    │
    ├── UserUsecase.CheckBalance
    │
    └── OrderUsecase.Create
```

Controller 不得关心这些操作如何实现，也不得知道：

```text
GoodsUsecase.CheckInventory → GoodsRepository.QueryGoodsInfo
UserUsecase.CheckBalance    → UserRepository.GetUserInfo
OrderUsecase.Create         → OrderRepository.Save
```

Controller 只依赖 Usecase 的契约。

### Controller 限制

Controller 不得：

* 直接访问数据库 / Redis / 消息队列 / 外部服务。
* 执行 SQL 或包含持久化逻辑。
* 包含 repository 实现细节。
* 实现模块内业务规则。
* 包含 HTTP/gRPC 特定逻辑。
* 依赖 Echo、net/http、gRPC 传输类型或其他协议特定类型。
* 知道 Usecase 是如何实现的。
* 重复实现 Usecase 的业务逻辑。

Controller 与外部协议完全解耦：协议转换由 Adapter 负责，Controller 的入参/出参必须是业务语义结构，同一工作流可被任意协议复用。详见 [adapter.md](adapter.md#controller-不感知外部协议)。

### Controller 复杂度

如果某个 Controller 因为包含复杂计算、数据库操作、详细业务规则、重复校验或模块内决策而变得庞大，应将那部分逻辑移入合适的 Usecase 或 Domain 抽象中。

Controller 读起来应当主要像一条业务工作流：

```text
校验输入
    ↓
检查库存
    ↓
检查余额
    ↓
创建订单
    ↓
返回结果
```

---

## Usecase

位置：`internal/app/usecase`

### 职责

Usecase 层包含**模块内业务逻辑的具体实现**。

例如：`GoodsUsecase`、`UserUsecase`、`OrderUsecase`、`PaymentUsecase`。

### Usecase 与 Repository

Usecase 可以依赖 Repository 接口/契约。

由 Usecase 决定：

* 需要哪些数据。
* 需要应用哪些业务规则。
* 需要执行哪些 repository 操作。
* 如何解读 repository 返回的数据。
* 业务操作成功还是失败。

示例流程：

```text
GoodsUsecase.CheckInventory
        │
        ▼
GoodsRepository.QueryGoodsInfo
        │
        ▼
检查 goods.inventory
        │
        ▼
返回结果 / ErrInsufficientInventory
```

示例：

```go
func (u *GoodsUsecase) CheckInventory(
	ctx context.Context,
	goodsID int64,
	quantity int,
) error {
	goods, err := u.goodsRepository.QueryGoodsInfo(ctx, goodsID)
	if err != nil {
		return fmt.Errorf("query goods info: %w", err)
	}

	if goods.Inventory < quantity {
		return errors.ErrInsufficientInventory
	}

	return nil
}
```

Repository 只负责取数据。数据的业务含义由 Usecase 决定。

### Usecase 限制

Usecase 不得：

* 依赖 HTTP 请求/响应对象、Echo context 或 gRPC 消息。
* 解析 HTTP 参数 / 构造 HTTP 响应 / 设置 HTTP 状态码。
* 实现路由或 cron 调度。
* 包含 SQL 实现。
* 直接访问 `database/sql`、Redis 客户端或消息队列实现。
* 包含传输层特定的 DTO 转换。

如果 Usecase 需要基础设施数据，必须使用 Repository 抽象。

### 模块归属

每个 Usecase 必须拥有其对应模块的业务规则。

```text
GoodsUsecase
    ├── CheckInventory
    ├── GetGoods
    └── CalculatePrice

UserUsecase
    ├── CheckBalance
    ├── GetUser
    └── ValidateUser

OrderUsecase
    ├── Create
    ├── Cancel
    └── Get
```

不要仅仅因为当前模块需要某个信息，就把业务逻辑移入另一个模块。

> 业务规则属于拥有该业务概念的模块。

---

## 跨模块操作

当某个操作涉及多个业务模块时，Controller 应当编排这些 Usecase。

不要为了简化 Controller 而创建这样的 Repository 方法：

```text
OrderRepository.CreateOrderAndCheckInventoryAndBalance
```

Repository 不是应用服务层。

---

## 禁止的捷径

### Adapter 直接查询 Repository

```text
handler → repository   // 禁止
```

> 特别规则：CRUD 例外仅适用于 controller → repository 这一步，adapter → repository 仍然一律禁止。

要求：

```text
handler → controller → usecase → repository
```

### Controller 直接查询数据库

```text
controller → gorm/db   // 禁止
```

> 特别规则：即使满足 CRUD 例外条件，Controller 也只能调用 Repository 抽象，不得直接操作数据库连接（gorm/db）。

### Controller 实现模块业务逻辑

```go
if goods.Inventory < quantity {
    ...
}
```

如果这是库存业务规则，它属于 `GoodsUsecase`，而不是 Controller。

### Usecase 访问 HTTP

```go
func (u *OrderUsecase) Create(c echo.Context, ...) // 禁止
func (u *OrderUsecase) Create(ctx context.Context, ...) // 要求
```

### Repository 实现业务工作流

```text
OrderRepository
    ↓ 检查库存
    ↓ 检查余额
    ↓ 创建订单
```

禁止。该工作流属于 Controller + 各 Usecase。
