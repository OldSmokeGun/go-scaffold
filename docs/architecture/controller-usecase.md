# Controller and Usecase

---

## Controller

Location: `internal/app/controller`

### Responsibility

The Controller is the **core application business entry point**.

It is responsible for:

1. Validating business input parameters.
2. Organizing business operations.
3. Orchestrating multiple Usecases.
4. Defining the execution order of Usecases.
5. Coordinating different business modules.
6. Controlling the application-level workflow.
7. Returning the final business result.

The Controller MUST NOT implement the internal business logic of individual modules.

* Controller defines **what business operations need to happen and in what order**.
* Usecase defines **how each individual business operation is implemented**.

### Controller MUST NOT Access Repository

Controllers MUST NOT directly call Repository, Database, Redis, Message Queue, or External Service.

Controllers MUST interact with business functionality through Usecases.

Forbidden:

```go
func (c *Controller) CreateOrder(...) error {
	goods, err := c.goodsRepository.QueryGoodsInfo(...)
	...
}
```

Required:

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

### Workflow Example

```text
CreateOrder
    │
    ├── GoodsUsecase.CheckInventory
    │
    ├── UserUsecase.CheckBalance
    │
    └── OrderUsecase.Create
```

The Controller MUST NOT care how these operations are implemented, and MUST NOT know:

```text
GoodsUsecase.CheckInventory → GoodsRepository.QueryGoodsInfo
UserUsecase.CheckBalance    → UserRepository.GetUserInfo
OrderUsecase.Create         → OrderRepository.Save
```

The Controller only depends on the Usecase contract.

### Controller Restrictions

The Controller MUST NOT:

* Access databases / Redis / message queues / external services directly.
* Execute SQL or contain persistence logic.
* Contain repository implementation details.
* Implement module-specific business rules.
* Contain HTTP/gRPC-specific logic.
* Depend on Echo, net/http, gRPC transport types, or other protocol-specific types.
* Know how a Usecase is implemented.
* Duplicate Usecase business logic.

### Controller Complexity

If a Controller becomes large because it contains complex calculations, database operations, detailed business rules, repeated validation, or module-specific decisions, move that logic into the appropriate Usecase or Domain abstraction.

A Controller should primarily read like a business workflow:

```text
validate input
    ↓
check inventory
    ↓
check balance
    ↓
create order
    ↓
return result
```

---

## Usecase

Location: `internal/app/usecase`

### Responsibility

The Usecase layer contains the **concrete implementation of module-specific business logic**.

Examples: `GoodsUsecase`, `UserUsecase`, `OrderUsecase`, `PaymentUsecase`.

### Usecase and Repository

A Usecase MAY depend on Repository interfaces/contracts.

The Usecase determines:

* What data is required.
* What business rules need to be applied.
* What repository operations need to be performed.
* How repository data is interpreted.
* Whether a business operation succeeds or fails.

Example flow:

```text
GoodsUsecase.CheckInventory
        │
        ▼
GoodsRepository.QueryGoodsInfo
        │
        ▼
Check goods.inventory
        │
        ▼
Return result / ErrInsufficientInventory
```

Example:

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

The Repository only retrieves data. The Usecase decides what the data means from a business perspective.

### Usecase Restrictions

The Usecase MUST NOT:

* Depend on HTTP request/response objects, Echo context, or gRPC messages.
* Parse HTTP parameters / construct HTTP responses / set HTTP status codes.
* Implement routing or cron scheduling.
* Contain SQL implementation.
* Directly access `database/sql`, Redis clients, or message queue implementations.
* Contain transport-specific DTO conversion.

If a Usecase requires infrastructure data, it MUST use a Repository abstraction.

### Module Ownership

Each Usecase MUST own the business rules of its corresponding module.

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

Do not move business logic into another module merely because the current module needs that information.

> A business rule belongs to the module that owns the business concept.

---

## Cross-Module Operations

When an operation involves multiple business modules, the Controller SHOULD orchestrate the Usecases.

Do NOT create a Repository method such as:

```text
OrderRepository.CreateOrderAndCheckInventoryAndBalance
```

merely to simplify the Controller. The Repository is not the application service layer.

---

## Prohibited Shortcuts

### Adapter directly querying Repository

```text
handler → repository   // forbidden
```

Required:

```text
handler → controller → usecase → repository
```

### Controller directly querying database

```text
controller → gorm/db   // forbidden
```

### Controller implementing module business logic

```go
if goods.Inventory < quantity {
    ...
}
```

If this is the inventory business rule, it belongs in `GoodsUsecase`, not Controller.

### Usecase accessing HTTP

```go
func (u *OrderUsecase) Create(c echo.Context, ...) // forbidden
func (u *OrderUsecase) Create(ctx context.Context, ...) // required
```

### Repository implementing business workflow

```text
OrderRepository
    ↓ check inventory
    ↓ check balance
    ↓ create order
```

Forbidden. This workflow belongs to Controller + Usecases.
