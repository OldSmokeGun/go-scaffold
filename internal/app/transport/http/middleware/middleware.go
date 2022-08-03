package middleware

// Logger middleware logger must implement this interface
type Logger interface {
	Debug(...any)
	Debugf(string, ...any)
	Info(...any)
	Infof(string, ...any)
	Error(...any)
	Errorf(string, ...any)
}

// ResponseBody middleware response body must implement this interface
type ResponseBody interface {
	WithMsg(msg string)
}
