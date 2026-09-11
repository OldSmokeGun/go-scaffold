# Change Workflow

---

## Before Modifying Code

1. Identify the responsible layer.
2. Check the dependency direction.
3. Read the relevant document under `docs/`.
4. Reuse existing project abstractions.
5. Make the smallest appropriate change.
6. Run formatting, linting and tests.

---

## New Feature Implementation Procedure

### Step 1 — Define the external entry point

Determine whether the feature is triggered by HTTP / gRPC / Cron / CLI / Script.

Implement the entry point in `adapter`.

### Step 2 — Define the Controller operation

Create the application-level business operation in `controller`.

The Controller should describe the workflow.

### Step 3 — Identify business modules

Determine which business modules participate, for example Goods / User / Order / Payment.

### Step 4 — Implement module business logic

Implement module-specific rules in the corresponding Usecases.

### Step 5 — Define data requirements

Determine what data each Usecase needs.

Expose the required Repository abstraction.

### Step 6 — Implement infrastructure access

Implement MySQL, Redis, MQ, or other infrastructure access inside Repository.

### Step 7 — Define domain concepts

If the feature introduces a stable business entity, enum, custom type, or business constant, define it in `domain`.

### Step 8 — Define business errors

If the feature introduces a new business error, define it in `internal/errors`.

---

## Mandatory Final Verification

Before considering a task complete, verify that:

1. Adapter only calls Controller.
2. Controller does not directly access Repository or infrastructure.
3. Controller only orchestrates Usecases.
4. Usecase contains module-specific business logic.
5. Usecase accesses data through Repository abstractions.
6. Repository contains infrastructure/data-access implementation only.
7. Domain remains independent from application infrastructure.
8. Business errors are centralized in `internal/errors`.
9. Transport-specific types do not leak into Controller or Usecase.
10. No new reverse dependency has been introduced.
11. No business workflow has been moved into Repository.
12. No business logic has been moved into Adapter.
13. No unnecessary business logic has been duplicated between Controller and Usecase.

Also run:

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
* Ignore architectural violations because they are inconvenient to fix.
