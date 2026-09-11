# Go Style and Development Rules

---

## Context

`context.Context` MAY be passed through the application layers:

```text
Adapter → Controller → Usecase → Repository
```

The Context MUST NOT be stored in long-lived structs.

Do not use:

```go
type OrderUsecase struct {
	ctx context.Context
}
```

Use:

```go
func (u *OrderUsecase) Create(
	ctx context.Context,
	input CreateOrderInput,
) (*domain.Order, error)
```

---

## Interfaces

* Keep interfaces small and meaningful.
* Consumers should depend on abstractions rather than concrete infrastructure types when appropriate.
* Prefer standard library solutions when appropriate.

---

## Functions and State

* Keep functions focused.
* Avoid unnecessary global mutable state.
* Avoid unnecessary `init()`.
* Avoid unrelated refactoring.
* Do not weaken architecture to make implementation easier.

---

## Errors

* Handle errors explicitly.
* Prefer `%w` when wrapping so the original error remains inspectable.
* Do not silently discard errors.
* See [../architecture/error-handling.md](../architecture/error-handling.md) for business error ownership.

---

## Code Quality

All code created or modified by the Agent MUST comply with the project's configured code quality and linting rules.

Before considering a task complete:

1. Run the project's configured lint and validation commands.
2. Fix all errors reported by the validation tools.
3. Do not disable, weaken, or bypass existing validation rules.
4. Do not modify lint configuration merely to make the code pass.
5. Do not use suppression directives unless there is a legitimate and documented reason.
6. Re-run validation after fixing issues.
7. Never claim that validation passed unless it was actually executed successfully.

For Go projects:

* Follow the project's `golangci-lint` configuration (`.golangci.yml`).
* Run at least:

```bash
golangci-lint fmt
golangci-lint run
```

* Fix all reported issues before completing the task.
