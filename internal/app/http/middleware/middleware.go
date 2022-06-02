package middleware

// Logger middleware logger must implement this interface
type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

// ResponseBody middleware response body must implement this interface
type ResponseBody interface {
	WithMsg(msg string)
}
