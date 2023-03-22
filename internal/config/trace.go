package config

// Trace the otel trace config
type Trace struct {
	Endpoint string `json:"endpoint"`
}

func (Trace) GetName() string {
	return "trace"
}
