# Go 风格与开发规则

---

## Context

`context.Context` 可以贯穿应用各层传递：

```text
Adapter → Controller → Usecase → Repository
```

Context 不得存储在长生命周期的结构体中。

不要使用：

```go
type OrderUsecase struct {
	ctx context.Context
}
```

应使用：

```go
func (u *OrderUsecase) Create(
	ctx context.Context,
	input CreateOrderInput,
) (*domain.Order, error)
```

---

## 接口

* 保持接口小而有意义。
* 在合适时，使用方应依赖抽象而不是具体的基础设施类型。
* 在合适时优先使用标准库方案。

---

## 函数与状态

* 保持函数职责单一。
* 避免不必要的全局可变状态。
* 避免不必要的 `init()`。
* 避免无关的重构。
* 不得为了让实现更容易而削弱架构。

---

## 错误

* 显式处理错误。
* 包装错误时优先使用 `%w`，使原始错误保持可检查。
* 不得静默丢弃错误。
* 业务错误归属参见 [../architecture/error-handling.md](../architecture/error-handling.md)。

---

## 代码质量

Agent 创建或修改的所有代码都必须符合项目配置的代码质量与 lint 规则。

在认定任务完成之前：

1. 运行项目配置的 lint 和验证命令。
2. 修复验证工具报告的所有错误。
3. 不得禁用、削弱或绕过既有验证规则。
4. 不得仅仅为了让代码通过而修改 lint 配置。
5. 除非有正当且已记录的理由，否则不得使用抑制指令。
6. 修复问题后重新运行验证。
7. 绝不在验证未实际执行成功的情况下声称验证已通过。

对于 Go 项目：

* 遵循项目的 `golangci-lint` 配置（`.golangci.yml`）。
* 至少运行：

```bash
golangci-lint fmt
golangci-lint run
```

* 在完成任务之前修复所有报告的问题。
