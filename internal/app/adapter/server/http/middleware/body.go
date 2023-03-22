package middleware

// Body formatted json response
type Body struct {
	ErrNo  int `json:"errNo,omitempty"`
	ErrMsg any `json:"errMsg,omitempty"`
	Stack  any `json:"stack,omitempty"`
	Data   any `json:"data,omitempty"`
}

// NewDefaultBody return json response body
func NewDefaultBody() *Body {
	return &Body{}
}

func (b *Body) WithErrNo(code int) *Body {
	b.ErrNo = code
	return b
}

func (b *Body) WithErrMsg(msg any) *Body {
	b.ErrMsg = msg
	return b
}

func (b *Body) WithStack(stack any) *Body {
	b.Stack = stack
	return b
}

func (b *Body) WithData(data any) *Body {
	b.Data = data
	return b
}
