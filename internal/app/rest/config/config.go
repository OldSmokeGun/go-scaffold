package config

type (
	Config struct {
		Host        string
		Port        int
		ExternalUrl string
		Jwt         Jwt
	}

	Jwt struct {
		Key string
	}
)
