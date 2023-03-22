package errors

// NoErrorCode code when there are no errors
const NoErrorCode = 10000

// standard errors
var (
	ErrServerError      = New("server error", 10001, "SERVER_ERROR")
	ErrBadRequest       = New("bad request", 10002, "BAD_REQUEST")
	ErrValidateError    = New("parameters validate error", 10003, "VALIDATE_ERROR")
	ErrUnauthorized     = New("unauthorized", 10004, "UNAUTHORIZED")
	ErrPermissionDenied = New("permission denied", 10005, "PERMISSION_DENIED")
	ErrResourceNotFound = New("resource not found", 10006, "RESOURCE_NOT_FOUND")
	ErrTooManyRequest   = New("too many request", 10007, "TOO_MANY_REQUEST")
)

// Error application internal error
type Error struct {
	msg   string
	code  int
	label string
}

// New returns an error that formats as the given text.
func New(text string, code int, label string) *Error {
	return &Error{text, code, label}
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Label() string {
	return e.label
}

func (e *Error) WithMsg(msg string) *Error {
	e.msg = msg
	return e
}
