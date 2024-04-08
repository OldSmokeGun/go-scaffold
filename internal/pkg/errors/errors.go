package errors

// standard errors
var (
	ErrInternalError      = New("internal error", 10001, "INTERNAL_ERROR")
	ErrBadCall            = New("bad call", 20001, "BAD_CALL")
	ErrValidateError      = New("parameters validate error", 20002, "VALIDATE_ERROR")
	ErrInvalidAuthorized  = New("invalid authorized", 20003, "INVALID_AUTHORIZED")
	ErrAccessDenied       = New("access denied", 20004, "ACCESS_DENIED")
	ErrResourceNotFound   = New("resource not found", 20005, "RESOURCE_NOT_FOUND")
	ErrCallsTooFrequently = New("call too frequently", 20006, "CALLS_TOO_FREQUENTLY")
)

// Error application internal error
type Error struct {
	msg   string
	code  int
	label string
	error error
}

// New returns an error that formats as the given text.
func New(text string, code int, label string) *Error {
	return &Error{text, code, label, nil}
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Cause() error {
	return e.error
}

func (e *Error) Unwrap() error {
	return e.error
}

func (e *Error) Msg() string {
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

func (e *Error) WithError(err error) *Error {
	e.error = err
	return e
}
