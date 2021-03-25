package config

import (
	"gin-scaffold/internal/config"
	"time"
)

type (
	Config struct {
		config.Config
		JwtConfig JwtConfig
	}

	JwtConfig struct {
		Key    string
		Expire time.Duration
	}
)
