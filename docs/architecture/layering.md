# 分层与依赖

这些规则是**强制性要求**，而非建议。

Agent 生成的或修改的每一段代码都必须符合此处定义的职责、依赖边界和数据流规则。

如果被要求的实现看起来与这些规则冲突，Agent 不得为了图方便而绕过它们。应当重新设计实现，使其遵守各层的职责。

---

## 目录结构

```text
internal
├─app
│  ├─adapter
│  │  ├─cron
│  │  │  ├─job
│  │  │  └─scheduler
│  │  ├─scripts
│  │  └─server
│  │      ├─grpc
│  │      │  ├─handler
│  │      │  └─router
│  │      └─http
│  │          ├─handler
│  │          ├─middleware
│  │          └─router
│  ├─controller
│  ├─domain
│  ├─repository
│  └─usecase
└─errors
```

分层流向：

```text
外部协议
   │
   ▼
Adapter（适配层）
   │
   ▼
Controller（控制层）
   │
   ▼
Usecase（用例层）
   │
   ▼
Repository（仓储层）
   │
   ▼
外部基础设施
（MySQL / Redis / MQ / HTTP 等）
```

* `Domain` 定义业务实体、值类型、枚举、常量以及稳定的业务概念。
* `errors` 定义应用/业务错误及其外部错误表示。

---

## 强制依赖方向

```text
adapter
   ↓
controller
   ↓
usecase
   ↓
repository
```

依赖方向不得逆转。

原则：

1. 下层不得依赖上层。
2. 基础设施实现不得包含应用编排。
3. 协议特定概念不得泄漏到业务逻辑中。
4. 业务逻辑不得依赖 HTTP、gRPC、CLI、cron 或其他传输协议。
5. Controller 不得直接访问数据库、Redis、消息队列或其他基础设施。
6. Adapter 不得实现业务逻辑。
7. Usecase 不得实现 HTTP/gRPC/CLI 响应处理。
8. Repository 实现不得执行应用级业务编排。
9. Domain 定义不得依赖 adapter、controller、usecase 或 repository。
10. 业务错误必须在 `internal/errors` 中集中定义。

---

## 依赖矩阵

| 来源       | 允许依赖的对象                        |
| ---------- | -------------------------------------------------- |
| adapter    | controller、domain、errors                         |
| controller | usecase、domain、errors                            |
| usecase    | repository、domain、errors                         |
| repository | domain、errors、基础设施库                         |
| domain     | 仅标准库，除非有明确正当理由                       |
| errors     | 仅标准库，除非有明确正当理由                       |

禁止的依赖：

| 禁止的依赖                                           | 理由                 |
| ----------------------- | ------------------------------ |
| adapter → usecase       | 绕过 Controller                |
| adapter → repository    | 绕过业务层                     |
| adapter → database      | 基础设施泄漏                   |
| controller → repository | 绕过 Usecase                   |
| controller → database   | 基础设施泄漏                   |
| controller → redis      | 基础设施泄漏                   |
| usecase → adapter       | 依赖逆转                       |
| usecase → HTTP / Echo   | 传输层泄漏                     |
| repository → controller | 依赖逆转                       |
| repository → usecase    | 业务编排泄漏                   |
| repository → adapter    | 依赖逆转                       |
| domain → repository / usecase / controller / adapter | Domain 被污染 |

---

## 业务逻辑归属

| 代码类型                                                       | 所属层     |
| --- | --- |
| 协议转换（HTTP / gRPC / CLI / Cron → 应用输入）                | adapter    |
| 业务工作流编排（`A → B → C → D`）                              | controller |
| 模块内业务规则（库存、余额、定价、权限）                       | usecase    |
| SQL / Redis / MQ / 基础设施 SDK 访问                           | repository |
| 业务实体 / 枚举 / 类型 / 常量                                  | domain     |
| 业务错误码 / 错误消息 / 状态映射                               | errors     |

---

## 架构决策顺序

当不确定代码属于哪一层时：

```text
它是否与特定协议相关？
        │
        ├── 是 → Adapter
        │
        └── 否
             │
             ▼
它是否编排多个业务操作？
        │
        ├── 是 → Controller
        │
        └── 否
             │
             ▼
它是否是模块内业务逻辑？
        │
        ├── 是 → Usecase
        │
        └── 否
             │
             ▼
它是否访问基础设施/数据？
        │
        ├── 是 → Repository
        │
        └── 否
             │
             ▼
它是否是稳定的业务实体/类型/值？
        │
        ├── 是 → Domain
        │
        └── 否
             │
             ▼
它是否是标准化的业务错误？
        │
        ├── 是 → errors
        │
        └── 否
             │
             ▼
在添加代码之前重新评估设计。
```

不得基于实现方便来选择所属层。

---

## 不可妥协的原则

> **Adapter 处理协议。Controller 编排业务操作。Usecase 实现模块业务逻辑。Repository 实现基础设施访问。Domain 定义业务概念。Errors 定义标准化业务错误。**

```text
                    ┌──────────────┐
                    │   Adapter    │
                    │ HTTP/gRPC/...│
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  Controller  │
                    │ Orchestration│
                    └──────┬───────┘
                           │
             ┌─────────────┼─────────────┐
             ▼             ▼             ▼
        ┌─────────┐   ┌─────────┐   ┌─────────┐
        │ Goods   │   │  User   │   │  Order  │
        │ Usecase │   │ Usecase │   │ Usecase │
        └────┬────┘   └────┬────┘   └────┬────┘
             │             │             │
             ▼             ▼             ▼
        ┌─────────┐   ┌─────────┐   ┌─────────┐
        │ Goods   │   │  User   │   │  Order  │
        │ Repo    │   │  Repo   │   │  Repo   │
        └────┬────┘   └────┬────┘   └────┬────┘
             │             │             │
             └─────────────┼─────────────┘
                           ▼
                    Infrastructure
                 MySQL / Redis / MQ / ...
```

任何违反此依赖方向的实现，除非本文档被明确修订，否则在架构上都是无效的。

---

## 重构规则

修改既有代码时：

1. 如果违规处于本次变更范围内，优先修复该违规。
2. 不得将既有违规扩散到新代码中。
3. 不得以既有违规为由制造新的违规。
4. 在需要时保持向后兼容，但新代码必须合规。

历史遗留的违规不构成制造新违规的许可。
