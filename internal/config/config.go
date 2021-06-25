package config

import (
	"time"
)

type (
	Config struct {
		AppConf AppConf
		JwtConf JwtConf
	}

	AppConf struct {
		Host string
		Port int
		Env  string
		Log  string
	}

	TemplateConf struct {
		Glob string
	}

	JwtConf struct {
		Key    string
		Expire time.Duration
	}
)
