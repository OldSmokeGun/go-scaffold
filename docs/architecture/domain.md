# Domain 层

位置：`internal/app/domain`

---

## 职责

Domain 层定义稳定的业务概念和业务实体。

它可以包含：

* 业务实体
* 值对象
* 自定义类型
* 枚举
* 业务常量
* 稳定的领域级结构
* 领域特定的基础类型

示例：

```go
type OrderStatus int

const (
	OrderStatusPending OrderStatus = 1
	OrderStatusPaid    OrderStatus = 2
	OrderStatusClosed  OrderStatus = 3
)
```

当业务常量代表稳定的业务规则或业务取值时，必须在 Domain 层中定义：

```go
const MinPayAmount = 1
```

---

## Domain 限制

Domain 层必须保持与应用基础设施无关。

Domain 代码不得依赖：

```text
HTTP
gRPC
Echo
MySQL
Redis
Repository 实现
Controller
Usecase
```

避免在 `domain` 中放置数据库特定模型、HTTP DTO 或传输层特定结构。
