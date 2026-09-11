# Domain Layer

Location: `internal/app/domain`

---

## Responsibility

The Domain layer defines stable business concepts and business entities.

It may contain:

* Business entities
* Value objects
* Custom types
* Enums
* Business constants
* Stable domain-level structures
* Domain-specific primitive types

Example:

```go
type OrderStatus int

const (
	OrderStatusPending OrderStatus = 1
	OrderStatusPaid    OrderStatus = 2
	OrderStatusClosed  OrderStatus = 3
)
```

Business constants MUST be defined in the Domain layer when they represent stable business rules or business values:

```go
const MinPayAmount = 1
```

---

## Domain Restrictions

The Domain layer MUST remain independent from application infrastructure.

Domain code MUST NOT depend on:

```text
HTTP
gRPC
Echo
MySQL
Redis
Repository implementations
Controller
Usecase
```

Avoid placing database-specific models, HTTP DTOs, or transport-specific structures in `domain`.
