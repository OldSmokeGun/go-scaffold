package uid

// Provide  snowflake generator
func Provide() *Uid {
	return New()
}
