# Adapter Layer

Location: `internal/app/adapter`

---

## Responsibility

The Adapter layer is the **external protocol conversion layer** and the external entry point of the application.

Adapters may include:

* HTTP
* gRPC
* cron
* CLI
* scripts
* other external invocation mechanisms

Adapters convert external representations into internal application representations, and convert internal results into external representations.

Example flow:

```text
HTTP request
    ↓
HTTP handler
    ↓
Controller request
```

```text
Controller result
    ↓
HTTP response
```

The Adapter layer MUST NOT contain business logic.

---

## Adapter MAY

* Parse HTTP request parameters / bodies.
* Parse gRPC request messages.
* Parse CLI arguments.
* Parse cron/job configuration.
* Perform protocol-specific validation.
* Convert external DTOs into Controller input structures.
* Convert Controller output into HTTP/gRPC/CLI response structures.
* Set HTTP status codes and headers.
* Convert business errors into protocol-specific error representations.
* Handle protocol-specific authentication/authorization middleware.
* Handle serialization / deserialization.
* Handle protocol-specific logging and tracing.

Valid example:

```go
func (h *Handler) CreateOrder(c echo.Context) error {
	var req CreateOrderRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	input := controller.CreateOrderInput{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	result, err := h.controller.CreateOrder(c.Request().Context(), input)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, toResponse(result))
}
```

---

## Adapter MUST NOT

* Query MySQL / Redis directly.
* Publish messages directly as part of business processing.
* Call Repository methods.
* Implement business rules or workflows.
* Calculate business prices / check inventory / check balance.
* Create business orders or modify business entities as business decisions.
* Orchestrate multiple Usecases.
* Decide business outcomes.

Mandatory call path:

```text
Adapter → Controller          // allowed
Adapter → Usecase             // forbidden
Adapter → Repository          // forbidden
Adapter → Database / Redis    // forbidden
```

An Adapter MUST NOT bypass the Controller even if doing so appears simpler.

---

## Adapter Subdirectories

```text
adapter
├─cron
│  ├─job
│  └─scheduler
├─scripts
└─server
   ├─grpc
   │  ├─handler
   │  └─router
   └─http
      ├─handler
      ├─middleware
      └─router
```

### `router`

* Route / endpoint registration.
* Protocol-level configuration.
* MUST NOT contain business logic.

### `handler`

* Receive external requests.
* Parse input.
* Call Controller.
* Convert Controller output to protocol responses.
* MUST NOT call Usecase or Repository directly.

### `middleware`

Protocol-level cross-cutting concerns:

* Authentication
* Request tracing
* Logging
* CORS
* Rate limiting
* Request metadata

Middleware MUST NOT implement business workflows.

### `cron/job`

Scheduled external entry points. MUST invoke Controller, not Usecase or Repository.

### `scheduler`

Scheduling and triggering jobs. MUST NOT contain business logic.

### `scripts`

External entry points. MUST follow the same Adapter rules and MUST NOT bypass Controller.

---

## DTO Rules

Transport-specific DTOs MUST remain inside Adapter, for example:

```text
adapter/http/handler
    CreateOrderRequest
    CreateOrderResponse
```

These types MUST NOT be used as business-layer entities.

Do not pass `*http.Request`, `echo.Context`, gRPC request/response types into Controller or Usecase.

Prefer internal request structures:

```go
type CreateOrderInput struct {
	UserID    int64
	ProductID int64
	Quantity  int
}
```

The Adapter converts external data into this structure.
