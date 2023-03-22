package config

// Discovery service discovery config
type Discovery struct {
	Etcd   *Etcd
	Consul *Consul
}

func (Discovery) GetName() string {
	return "discovery"
}

// Etcd the etcd config
type Etcd struct {
	Endpoints []string `json:"endpoints"`
}

// Consul the consul config
type Consul struct {
	Addr   string `json:"addr"`
	Schema string `json:"schema"`
}
