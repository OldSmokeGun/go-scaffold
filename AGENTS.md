# 项目架构

所有代码必须遵循分层架构：

```text
adapter → controller → usecase → repository
```

职责划分：

* adapter：协议转换与外部入口。
* controller：业务工作流编排。
* usecase：模块内具体业务逻辑。
* repository：基础设施与数据访问。
* domain：业务实体、类型、枚举与常量。
* errors：标准化业务错误。

各层详细职责与依赖规则：

> 参见 [docs/architecture/layering.md](docs/architecture/layering.md)

---

# Controller 与 Usecase

* controller 必须编排业务工作流。
* usecase 必须实现模块内具体业务逻辑。
* controller 不得直接访问 repository 或基础设施。
* usecase 不得依赖 adapter 或传输层特定类型。
* 跨模块工作流必须由 controller 编排。
* 例外：极简单的单资源 CRUD 可由 controller 直接调用 repository，避免空壳 Usecase 的样板代码（详见 controller-usecase.md）。

详细规则与示例：

> 参见 [docs/architecture/controller-usecase.md](docs/architecture/controller-usecase.md)

---

# Repository

* repository 必须处理数据访问与基础设施集成。
* repository 不得实现业务工作流。
* repository 不得依赖 controller 或 usecase。
* 数据库、Redis、MQ 及外部服务访问必须保留在 repository 边界之内。

详细规则：

> 参见 [docs/architecture/repository.md](docs/architecture/repository.md)

---

# Domain

* domain 包含业务实体、类型、枚举与常量。
* domain 不得依赖应用层或基础设施层。
* domain 应保持与传输层和持久化实现无关。

详细规则：

> 参见 [docs/architecture/domain.md](docs/architecture/domain.md)

---

# 错误处理

* 业务错误必须在 `internal/errors` 中定义。
* 业务错误不得包含业务工作流。
* 基础设施错误应被包装并在传播时保留原始错误。
* 传输层特定的错误呈现属于 adapter。

详细规则：

> 参见 [docs/architecture/error-handling.md](docs/architecture/error-handling.md)

---

# Adapter

* adapter 是外部协议入口（HTTP / gRPC / cron / CLI / 脚本）。
* adapter 只负责外部协议与 Controller 入参/出参之间的双向转换，不得编排业务逻辑。
* controller 不得感知外部协议，入参/出参必须是业务语义结构，同一工作流可被任意协议复用。
* adapter 只能调用 controller。
* adapter 不得直接调用 usecase、repository 或基础设施。
* 传输层特定的 DTO 必须保留在 adapter 内部。

详细规则与示例：

> 参见 [docs/architecture/adapter.md](docs/architecture/adapter.md)

---

# 通用开发规则

所有 Go 代码必须符合惯用写法、简洁且易于维护。

* 显式处理错误。
* 正确传递 `context.Context`。
* 保持函数职责单一。
* 保持接口小而有意义。
* 避免不必要的全局可变状态。
* 避免不必要的 `init()`。
* 在合适时优先使用标准库方案。
* 避免无关的重构。
* 不得为了让实现更容易而削弱架构。
* 不得仅仅为了让检查通过而禁用 lint 规则。

Go 开发详细规则：

> 参见 [docs/development/go-style.md](docs/development/go-style.md)

---

# 代码质量

项目使用 `.golangci.yml` 进行 Go 格式化和静态分析。

在完成非简单任务之前，运行：

```bash
golangci-lint fmt
golangci-lint run
go test ./...
```

所有 lint 错误必须解决。

禁止事项：

* 仅仅为了让检查通过而禁用某个 linter。
* 在没有正当且已记录理由的情况下添加 `//nolint`。
* 仅仅为了绕过既有违规而修改 `.golangci.yml`。
* 除非验证确实执行成功，否则不得声称验证已通过。

Lint 与格式化配置：

> 参见 [.golangci.yml](.golangci.yml)

---

# 变更工作流

修改代码之前：

1. 确定所属的职责层。
2. 检查依赖方向。
3. 阅读 `docs/architecture/` 或 `docs/development/` 下的相关文档。
4. 复用项目已有的抽象。
5. 做出最小的适当变更。
6. 运行格式化、lint 和测试。

新功能实现流程与最终验证清单：

> 参见 [docs/development/change-workflow.md](docs/development/change-workflow.md)

---

# 权威来源

`AGENTS.md` 定义强制性的高层规则。

详细规则定义在：

```text
docs/architecture/
docs/development/
```

自动化代码质量规则定义在：

```text
.golangci.yml
```

当需要某条详细规则时，阅读被引用的文档，而不是凭假设行事。
