package config

// JWT the JWT config
type JWT struct {
	Key string `json:"key"`
}

func (JWT) GetName() string {
	return "jwt"
}
