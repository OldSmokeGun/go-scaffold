# 变更工作流

---

## 修改代码之前

1. 确定所属的职责层。
2. 检查依赖方向。
3. 阅读 `docs/` 下的相关文档。
4. 复用项目已有的抽象。
5. 做出最小的适当变更。
6. 运行格式化、lint 和测试。

---

## 新功能实现流程

### 步骤 1 — 确定外部入口

确定该功能由 HTTP / gRPC / Cron / CLI / 脚本 中的哪种方式触发。

在 `adapter` 中实现入口。

### 步骤 2 — 定义 Controller 操作

在 `controller` 中创建应用级业务操作。

Controller 应当描述工作流。

### 步骤 3 — 确定业务模块

确定哪些业务模块参与，例如 商品 / 用户 / 订单 / 支付。

### 步骤 4 — 实现模块业务逻辑

在相应的 Usecase 中实现模块内业务规则。

### 步骤 5 — 确定数据需求

确定每个 Usecase 需要什么数据。

暴露所需的 Repository 抽象。

### 步骤 6 — 实现基础设施访问

在 Repository 内实现 MySQL、Redis、MQ 或其他基础设施访问。

### 步骤 7 — 定义领域概念

如果该功能引入了稳定的业务实体、枚举、自定义类型或业务常量，在 `domain` 中定义。

### 步骤 8 — 定义业务错误

如果该功能引入了新的业务错误，在 `internal/errors` 中定义。

---

## 强制性最终验证

在认定任务完成之前，确认：

1. Adapter 只调用 Controller。
2. Controller 不直接访问 Repository 或基础设施。
3. Controller 只编排 Usecase。
4. Usecase 包含模块内业务逻辑。
5. Usecase 通过 Repository 抽象访问数据。
6. Repository 只包含基础设施/数据访问实现。
7. Domain 保持与应用基础设施无关。
8. 业务错误集中在 `internal/errors` 中。
9. 传输层特定类型未泄漏到 Controller 或 Usecase。
10. 未引入新的逆向依赖。
11. 没有业务工作流被移入 Repository。
12. 没有业务逻辑被移入 Adapter。
13. Controller 与 Usecase 之间没有重复不必要的业务逻辑。

同时运行：

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
* 因为修复架构违规不方便而忽略它们。
