# Layering and Dependencies

These rules are **strict requirements**, not recommendations.

Every Agent-generated or Agent-modified code MUST comply with the responsibilities, dependency boundaries, and data-flow rules defined here.

If a requested implementation appears to conflict with these rules, the Agent MUST NOT bypass them for convenience. Redesign the implementation to respect layer responsibilities.

---

## Directory Structure

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

Layer flow:

```text
External Protocol
       │
       ▼
   Adapter
       │
       ▼
   Controller
       │
       ▼
    Usecase
       │
       ▼
   Repository
       │
       ▼
External Infrastructure
(MySQL / Redis / MQ / HTTP / etc.)
```

* `Domain` defines business entities, value types, enums, constants, and stable business concepts.
* `errors` defines application/business errors and their external error representation.

---

## Mandatory Dependency Direction

```text
adapter
   ↓
controller
   ↓
usecase
   ↓
repository
```

The dependency direction MUST NOT be reversed.

Principles:

1. Lower layers MUST NOT depend on higher layers.
2. Infrastructure implementations MUST NOT contain application orchestration.
3. Protocol-specific concepts MUST NOT leak into business logic.
4. Business logic MUST NOT depend on HTTP, gRPC, CLI, cron, or other transport protocols.
5. Controllers MUST NOT directly access databases, Redis, message queues, or other infrastructure.
6. Adapters MUST NOT implement business logic.
7. Usecases MUST NOT implement HTTP/gRPC/CLI response handling.
8. Repository implementations MUST NOT perform application-level business orchestration.
9. Domain definitions MUST NOT depend on adapters, controllers, usecases, or repositories.
10. Business errors MUST be defined centrally in `internal/errors`.

---

## Dependency Matrix

| From       | Allowed to depend on                               |
| ---------- | -------------------------------------------------- |
| adapter    | controller, domain, errors                         |
| controller | usecase, domain, errors                            |
| usecase    | repository, domain, errors                         |
| repository | domain, errors, infrastructure libraries           |
| domain     | standard library only, unless explicitly justified |
| errors     | standard library only, unless explicitly justified |

Forbidden:

| Forbidden dependency    | Reason                         |
| ----------------------- | ------------------------------ |
| adapter → usecase       | Bypasses Controller            |
| adapter → repository    | Bypasses business layer        |
| adapter → database      | Infrastructure leakage         |
| controller → repository | Bypasses Usecase               |
| controller → database   | Infrastructure leakage         |
| controller → redis      | Infrastructure leakage         |
| usecase → adapter       | Reversed dependency            |
| usecase → HTTP / Echo   | Transport leakage              |
| repository → controller | Reversed dependency            |
| repository → usecase    | Business orchestration leakage |
| repository → adapter    | Reversed dependency            |
| domain → repository / usecase / controller / adapter | Domain contamination |

---

## Business Logic Placement

| Kind of code | Layer |
| --- | --- |
| Protocol conversion (HTTP / gRPC / CLI / Cron → app input) | adapter |
| Business workflow orchestration (`A → B → C → D`) | controller |
| Module-specific business rules (inventory, balance, pricing, permission) | usecase |
| SQL / Redis / MQ / infrastructure SDK access | repository |
| Business entity / enum / type / constant | domain |
| Business error code / message / status mapping | errors |

---

## Architectural Decision Order

When uncertain where code belongs:

```text
Is it protocol-specific?
        │
        ├── YES → Adapter
        │
        └── NO
             │
             ▼
Does it orchestrate multiple business operations?
        │
        ├── YES → Controller
        │
        └── NO
             │
             ▼
Is it module-specific business logic?
        │
        ├── YES → Usecase
        │
        └── NO
             │
             ▼
Does it access infrastructure/data?
        │
        ├── YES → Repository
        │
        └── NO
             │
             ▼
Is it a stable business entity/type/value?
        │
        ├── YES → Domain
        │
        └── NO
             │
             ▼
Is it a standardized business error?
        │
        ├── YES → errors
        │
        └── NO
             │
             ▼
Re-evaluate the design before adding the code.
```

Do NOT choose a layer based on implementation convenience.

---

## Non-Negotiable Principle

> **Adapter handles protocols. Controller orchestrates business operations. Usecase implements module business logic. Repository implements infrastructure access. Domain defines business concepts. Errors defines standardized business errors.**

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

Any implementation that violates this dependency direction is architecturally invalid unless this document is explicitly amended.

---

## Refactoring Rules

When modifying existing code:

1. Prefer fixing a violation if it is within the scope of the requested change.
2. Do not spread an existing violation into new code.
3. Do not use an existing violation as justification for creating another one.
4. Preserve backward compatibility when required, but keep new code compliant.

A legacy violation does not create permission for new violations.
