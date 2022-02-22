package config

type Remote struct {
	Type          string // 远程配置类型
	Endpoint      string // 远程配置地址
	Path          string // 远程配置 path
	SecretKeyring string // 密钥
	Options       map[string]interface{}
}
