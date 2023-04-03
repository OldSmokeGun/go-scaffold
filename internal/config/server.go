package config

import "time"

// HTTP the HTTP config
type HTTP struct {
	Server *HTTPServer `json:"server"`
	JWT    *JWT        `json:"jwt"`
	Casbin *Casbin     `json:"casbin"`
}

func (HTTP) GetName() string {
	return "http"
}

// HTTPServer the HTTP server config
type HTTPServer struct {
	Network      string        `json:"network"`
	Addr         string        `json:"addr"`
	Timeout      time.Duration `json:"timeout"`
	ExternalAddr string        `json:"externalAddr"`
}

func (HTTPServer) GetName() string {
	return "http.server"
}

// JWT the JWT config
type JWT struct {
	Key string `json:"key"`
}

func (JWT) GetName() string {
	return "http.jwt"
}

// Casbin casbin config
type Casbin struct {
	Model   CasbinModel   `json:"model"`
	Adapter CasbinAdapter `json:"adapter"`
}

func (Casbin) GetName() string {
	return "http.casbin"
}

// CasbinModel casbin model
type CasbinModel struct {
	Path string `json:"path"`
}

// CasbinFileAdapter casbin file adapter
type (
	// CasbinAdapter casbin adapter
	CasbinAdapter struct {
		File *CasbinFileAdapter `json:"file"`
		Gorm *CasbinGormAdapter `json:"gorm"`
	}

	CasbinFileAdapter struct {
		Path string `json:"path"`
	}

	// CasbinGormAdapter casbin gorm adapter
	CasbinGormAdapter struct {
		TableName string `json:"tableName"`
	}
)

// GRPC the gRPC config
type GRPC struct {
	Server *GRPCServer `json:"server"`
}

func (GRPC) GetName() string {
	return "grpc"
}

// GRPCServer the gRPC server config
type GRPCServer struct {
	Network string        `json:"network"`
	Addr    string        `json:"addr"`
	Timeout time.Duration `json:"timeout"`
}

func (GRPCServer) GetName() string {
	return "grpc.server"
}
