package remote

import (
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/spf13/viper"
	"go-scaffold/pkg/config/remote/apollo"
	"strings"
)

type Remote struct {
	Type          string // 远程配置类型
	Endpoint      string // 远程配置地址
	Path          string // 配置 path
	SecretKeyring string // 密钥
	Options       map[string]interface{}
	ConfigType    string // 配置文件类型

	viper *viper.Viper
}

func New(t, endpoint, path, sk, configType string, ops ...OptionFunc) (*Remote, error) {
	r := &Remote{
		Type:          t,
		Endpoint:      endpoint,
		Path:          path,
		SecretKeyring: sk,
		ConfigType:    configType,
		viper:         viper.New(),
	}

	for _, op := range ops {
		op(r)
	}

	if r.Type == "apollo" {
		if err := initApolloProvider(r.Options); err != nil {
			return nil, err
		}
		apollo.WithConfigType(r.ConfigType)
	}

	return r, nil
}

type OptionFunc func(remote *Remote)

func WithOptions(options map[string]interface{}) OptionFunc {
	return func(remote *Remote) {
		remote.Options = options
	}
}

// Load 加载配置到模型
func (r *Remote) Load(model interface{}) error {
	if r.SecretKeyring == "" {
		if err := r.viper.AddRemoteProvider(r.Type, r.Endpoint, r.Path); err != nil {
			return err
		}
	} else {
		if err := r.viper.AddSecureRemoteProvider(r.Type, r.Endpoint, r.Path, r.SecretKeyring); err != nil {
			return err
		}
	}

	r.viper.SetConfigType(r.ConfigType)

	if err := r.viper.ReadRemoteConfig(); err != nil {
		return err
	}

	if err := r.viper.Unmarshal(&model); err != nil {
		return err
	}

	return nil
}

// MustLoad 加载配置到模型
func (r *Remote) MustLoad(model interface{}) {
	if err := r.Load(model); err != nil {
		panic(err)
	}
}

type OnConfigChange func(*viper.Viper, interface{})

// Watch 监控远程配置变更
func (r *Remote) Watch(f ...OnConfigChange) error {
	if r.Type == "apollo" {
		if len(f) > 0 {
			if f[0] != nil {
				apollo.WithOnConfigChange(func(event *storage.ChangeEvent) {
					f[0](r.viper, event)
				})
			}
		}

		if len(f) > 1 {
			if f[1] != nil {
				apollo.WithOnNewestConfigChange(func(event *storage.FullChangeEvent) {
					f[1](r.viper, event)
				})
			}
		}
	}

	if err := r.viper.WatchRemoteConfigOnChannel(); err != nil {
		return err
	}

	return nil
}

// MustWatch 监控远程配置变更
func (r *Remote) MustWatch(f ...OnConfigChange) {
	if err := r.Watch(f...); err != nil {
		panic(err)
	}
}

func initApolloProvider(options map[string]interface{}) error {
	appID, ok := options[strings.ToLower("AppID")]
	if !ok {
		return FieldMissingError{"AppID"}
	}

	appIDString, ok := appID.(string)
	if !ok {
		return FieldTypeConvertError{"AppID", "string"}
	}
	apollo.WithAppID(appIDString)

	cluster, ok := options[strings.ToLower("Cluster")]
	if !ok {
		return FieldMissingError{"Cluster"}
	}

	clusterString, ok := cluster.(string)
	if !ok {
		return FieldTypeConvertError{"Cluster", "string"}
	}
	if cluster != "" {
		apollo.WithCluster(clusterString)
	}

	return nil
}
