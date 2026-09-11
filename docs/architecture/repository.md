# Repository Layer

Location: `internal/app/repository`

---

## Responsibility

The Repository layer is responsible for **data access and infrastructure implementation**.

Implementations may include:

* MySQL / PostgreSQL
* Redis
* Message queue producers / consumers
* External data sources
* Other persistence/infrastructure mechanisms

Examples: `GoodsRepository`, `UserRepository`, `OrderRepository`.

The Repository translates application-level data access requirements into infrastructure-specific operations.

---

## Repository MAY

* Execute SQL.
* Query / update MySQL or Redis.
* Publish / consume messages when the repository abstraction is appropriate.
* Call infrastructure SDKs.
* Convert database models into application/domain representations.
* Handle infrastructure-specific errors.
* Manage persistence-specific details.

Example:

```go
func (r *GoodsRepository) QueryGoodsInfo(
	ctx context.Context,
	goodsID int64,
) (*domain.Goods, error) {
	...
}
```

The Usecase should not need to know whether this operation uses MySQL, Redis, Cache + MySQL, HTTP, or RPC.

---

## Repository MUST NOT

* Contain Controller workflows.
* Call Controllers.
* Call unrelated Usecases to perform business orchestration.
* Decide application-level business outcomes.
* Perform HTTP / gRPC response handling.
* Know about Echo or HTTP status codes.
* Return transport-specific errors.
* Implement business workflows spanning multiple modules.

Forbidden example:

```go
func (r *OrderRepository) CreateOrder(...) error {
	// check user balance
	// check inventory
	// create order
	// deduct balance
}
```

Those operations belong to the appropriate Usecases and Controller orchestration.

---

## Interface and Implementation

Interfaces MUST be defined according to the dependency direction.

A consumer should depend on an abstraction rather than an infrastructure implementation:

```text
Usecase
   ↓
GoodsRepository interface
   ↓
MySQL implementation
```

The Usecase MUST NOT depend directly on a concrete MySQL repository implementation when an abstraction is appropriate.

Infrastructure-specific implementation details MUST remain inside Repository.
