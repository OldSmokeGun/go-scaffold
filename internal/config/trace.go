package config

// Trace the trace config
type Trace struct {
	Protocol string `json:"protocol"`
	Endpoint string `json:"endpoint"`
}

func (Trace) GetName() string {
	return "trace"
}
