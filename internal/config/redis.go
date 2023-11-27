package config

import "time"

// Redis is redis config
type Redis struct {
	Addr               string        `json:"addr"`
	Username           string        `json:"username"`
	Password           string        `json:"password"`
	Database           int           `json:"database"`
	MaxRetries         int           `json:"maxRetries"`
	MinRetryBackoff    time.Duration `json:"minRetryBackoff"`
	MaxRetryBackoff    time.Duration `json:"maxRetryBackoff"`
	DialTimeout        time.Duration `json:"dialTimeout"`
	ReadTimeout        time.Duration `json:"readTimeout"`
	WriteTimeout       time.Duration `json:"writeTimeout"`
	PoolSize           int           `json:"poolSize"`
	MinIdleConns       int           `json:"minIdleConns"`
	MaxConnAge         time.Duration `json:"maxConnAge"`
	PoolTimeout        time.Duration `json:"poolTimeout"`
	IdleTimeout        time.Duration `json:"idleTimeout"`
	IdleCheckFrequency time.Duration `json:"idleCheckFrequency"`
}

func (Redis) GetName() string {
	return "redis"
}
