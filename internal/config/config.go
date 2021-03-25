package config

type (
	Config struct {
		ServerConfig
	}

	ServerConfig struct {
		Host         string
		Port         int
		Env          string
		Log          string
		TemplateGlob string
	}
)
