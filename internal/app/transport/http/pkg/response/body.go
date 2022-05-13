package response

// BodyInterface 响应需实现的接口
type BodyInterface interface {
	WithCode(code int)
	WithMsg(msg string)
	WithData(data interface{})
}

// Body 响应格式
type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// NewBody 返回 BodyInterface
func NewBody(code int, msg string, data interface{}) *Body {
	return &Body{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// WithCode 设置 Body 的 Code
func (b *Body) WithCode(code int) {
	b.Code = code
}

// WithMsg 设置 Body 的 Msg
func (b *Body) WithMsg(msg string) {
	b.Msg = msg
}

// WithData 设置 Body 的 Data
func (b *Body) WithData(data interface{}) {
	b.Data = data
}

// Option Body 属性设置函数
type Option func(*Body)

// WithCode 设置 Body 的 Code
func WithCode(code int) Option {
	return func(p *Body) {
		p.Code = code
	}
}

// WithMsg 设置 Body 的 Msg
func WithMsg(msg string) Option {
	return func(p *Body) {
		p.Msg = msg
	}
}

// WithData 设置 Body 的 Data
func WithData(data interface{}) Option {
	return func(p *Body) {
		p.Data = data
	}
}
