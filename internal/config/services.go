package config

// Services service list
type Services struct {
	Self string `json:"self"`
}

func (Services) GetName() string {
	return "services"
}
