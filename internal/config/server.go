package config

import "time"

// HTTP the HTTP server config
type HTTP struct {
	Network      string        `json:"network"`
	Addr         string        `json:"addr"`
	Timeout      time.Duration `json:"timeout"`
	ExternalAddr string        `json:"externalAddr"`
}

func (HTTP) GetName() string {
	return "http"
}

// GRPC the gRPC server config
type GRPC struct {
	Network string        `json:"network"`
	Addr    string        `json:"addr"`
	Timeout time.Duration `json:"timeout"`
}

func (GRPC) GetName() string {
	return "grpc"
}
