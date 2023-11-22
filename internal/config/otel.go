package config

// OTLP the otpl config
type OTLP struct {
	Trace Trace `json:"trace"`
}

func (OTLP) GetName() string {
	return "otpl"
}

// Trace the otpl trace config
type Trace struct {
	Protocol string `json:"protocol"`
	Endpoint string `json:"endpoint"`
}
