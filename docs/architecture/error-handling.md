# Error Handling

Location: `internal/errors`

---

## Responsibility

The `errors` package defines standardized application and business errors.

It may contain:

* Business error codes
* Business error messages
* HTTP status codes associated with business errors
* Standard application errors
* Error constructors
* Predefined business errors

Example:

```go
var ErrInsufficientBalance = BizError(
	40001,
	400,
	"余额不足",
)
```

---

## Error Ownership

Business errors MUST be defined centrally.

Do NOT recreate the same business error in multiple Usecases, Controllers, or Adapters.

Forbidden:

```go
return errors.New("余额不足")
```

when a corresponding centralized business error already exists.

Required:

```go
return errors.ErrInsufficientBalance
```

---

## Error Layer Restrictions

The `errors` package MUST NOT:

* Access repositories or databases.
* Call Usecases or Controllers.
* Handle HTTP requests.
* Execute business workflows.
* Depend on transport implementations.

The Adapter may translate a business error into a protocol-specific representation:

```text
internal/errors.ErrInsufficientBalance
              │
              ▼
HTTP Adapter
              │
              ▼
HTTP 400 + response body
```

The HTTP status code mapping MUST NOT be implemented inside the Usecase.

---

## Error Propagation

Errors MUST preserve their original cause whenever additional context is added.

Use:

```go
return fmt.Errorf("query goods info: %w", err)
```

Do NOT use:

```go
return fmt.Errorf("query goods info: %v", err)
```

when the original error needs to remain inspectable.

Do NOT silently discard errors.

* Business errors defined in `internal/errors` SHOULD be returned directly when no additional context is necessary.
* Infrastructure errors SHOULD normally be wrapped with meaningful context before being returned to the caller.
