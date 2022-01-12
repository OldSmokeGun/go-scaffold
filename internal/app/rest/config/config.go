package config

import (
	"time"
)

type (
	Config struct {
		Host        string
		Port        int
		ExternalUrl string
		Jwt         Jwt
	}

	Jwt struct {
		Key    string
		Expire time.Duration
	}
)
