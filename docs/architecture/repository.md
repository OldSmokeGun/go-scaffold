# Repository 层

位置：`internal/app/repository`

---

## 职责

Repository 层负责**数据访问与基础设施实现**。

实现可能包括：

* MySQL / PostgreSQL
* Redis
* 消息队列生产者 / 消费者
* 外部数据源
* 其他持久化/基础设施机制

例如：`GoodsRepository`、`UserRepository`、`OrderRepository`。

Repository 将应用级的数据访问需求转化为基础设施特定的操作。

---

## Repository 可以做的事

* 执行 SQL。
* 查询 / 更新 MySQL 或 Redis。
* 在仓储抽象合适的情况下发布 / 消费消息。
* 调用基础设施 SDK。
* 将数据库模型转换为应用/domain 表示。
* 处理基础设施特定的错误。
* 管理持久化相关的细节。

示例：

```go
func (r *GoodsRepository) QueryGoodsInfo(
	ctx context.Context,
	goodsID int64,
) (*domain.Goods, error) {
	...
}
```

Usecase 无需知道该操作使用的是 MySQL、Redis、Cache + MySQL、HTTP 还是 RPC。

---

## Repository 不得做的事

* 包含 Controller 工作流。
* 调用 Controller。
* 调用无关的 Usecase 来执行业务编排。
* 决定应用级业务结果。
* 执行 HTTP / gRPC 响应处理。
* 知道 Echo 或 HTTP 状态码。
* 返回传输层特定的错误。
* 实现跨多个模块的业务工作流。

禁止的示例：

```go
func (r *OrderRepository) CreateOrder(...) error {
	// 检查用户余额
	// 检查库存
	// 创建订单
	// 扣减余额
}
```

这些操作属于相应的 Usecase 和 Controller 编排。

---

## 接口与实现

必须按照依赖方向来定义接口。

使用方应当依赖抽象，而不是基础设施实现：

```text
Usecase
   ↓
GoodsRepository 接口
   ↓
MySQL 实现
```

在抽象合适的情况下，Usecase 不得直接依赖具体的 MySQL repository 实现。

基础设施特定的实现细节必须保留在 Repository 内部。
