# Project Architecture

All code MUST follow the layered architecture:

```text
adapter → controller → usecase → repository
```

Responsibilities:

* adapter: protocol conversion and external entry points.
* controller: business workflow orchestration.
* usecase: module-specific business logic.
* repository: infrastructure and data access.
* domain: business entities, types, enums and constants.
* errors: standardized business errors.

Detailed layer responsibilities and dependency rules:

> See [docs/architecture/layering.md](docs/architecture/layering.md)

---

# Controller and Usecase

* controller MUST orchestrate business workflows.
* usecase MUST implement module-specific business logic.
* controller MUST NOT access repository or infrastructure directly.
* usecase MUST NOT depend on adapter or transport-specific types.
* Cross-module workflows MUST be orchestrated by controller.

Detailed rules and examples:

> See [docs/architecture/controller-usecase.md](docs/architecture/controller-usecase.md)

---

# Repository

* repository MUST handle data access and infrastructure integration.
* repository MUST NOT implement business workflows.
* repository MUST NOT depend on controller or usecase.
* Database, Redis, MQ and external service access MUST remain behind repository boundaries.

Detailed rules:

> See [docs/architecture/repository.md](docs/architecture/repository.md)

---

# Domain

* domain contains business entities, types, enums and constants.
* domain MUST NOT depend on application or infrastructure layers.
* domain SHOULD remain independent from transport and persistence implementations.

Detailed rules:

> See [docs/architecture/domain.md](docs/architecture/domain.md)

---

# Errors

* Business errors MUST be defined in `internal/errors`.
* Business errors MUST NOT contain business workflows.
* Infrastructure errors SHOULD be wrapped and propagated without losing the original error.
* Transport-specific error presentation belongs to adapter.

Detailed rules:

> See [docs/architecture/error-handling.md](docs/architecture/error-handling.md)

---

# Adapter

* adapter is the external protocol entry point (HTTP / gRPC / cron / CLI / scripts).
* adapter MUST only call controller.
* adapter MUST NOT call usecase, repository, or infrastructure directly.
* Transport-specific DTOs MUST remain inside adapter.

Detailed rules and examples:

> See [docs/architecture/adapter.md](docs/architecture/adapter.md)

---

# General Development Rules

All Go code MUST be idiomatic, simple and maintainable.

* Handle errors explicitly.
* Propagate `context.Context` correctly.
* Keep functions focused.
* Keep interfaces small and meaningful.
* Avoid unnecessary global mutable state.
* Avoid unnecessary `init()`.
* Prefer standard library solutions when appropriate.
* Avoid unrelated refactoring.
* Do not weaken architecture to make implementation easier.
* Do not disable lint rules merely to make code pass.

Detailed Go development rules:

> See [docs/development/go-style.md](docs/development/go-style.md)

---

# Code Quality

The project uses `.golangci.yml` for Go formatting and static analysis.

Before completing a non-trivial task, run:

```bash
golangci-lint fmt
golangci-lint run
go test ./...
```

All lint errors MUST be resolved.

Do NOT:

* Disable a linter merely to make the check pass.
* Add `//nolint` without a legitimate and documented reason.
* Modify `.golangci.yml` merely to bypass an existing violation.
* Claim validation passed unless it was actually executed successfully.

Lint and formatting configuration:

> See [.golangci.yml](.golangci.yml)

---

# Change Workflow

Before modifying code:

1. Identify the responsible layer.
2. Check the dependency direction.
3. Read the relevant document under `docs/architecture/` or `docs/development/`.
4. Reuse existing project abstractions.
5. Make the smallest appropriate change.
6. Run formatting, linting and tests.

New feature implementation procedure and final verification checklist:

> See [docs/development/change-workflow.md](docs/development/change-workflow.md)

---

# Source of Truth

`AGENTS.md` defines mandatory high-level rules.

Detailed rules are defined in:

```text
docs/architecture/
docs/development/
```

Automated code quality rules are defined in:

```text
.golangci.yml
```

When a detailed rule is needed, read the referenced document instead of making assumptions.
