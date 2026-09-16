package uid

// Provide  sonyflake generator
func Provide() (*UID, error) {
	return New()
}
