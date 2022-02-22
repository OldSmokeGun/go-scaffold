package apollo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/cluster"
	"github.com/apolloconfig/agollo/v4/component/log"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/env/file"
	"github.com/apolloconfig/agollo/v4/protocol/auth"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/magiconair/properties"
	crypt "github.com/sagikazarmark/crypt/config"
	"gopkg.in/yaml.v2"
)

var (
	ErrApolloAppIDMissing = errors.New("the apollo appID missing")
)

type ConfigType string

const (
	Properties ConfigType = "properties"
	Json       ConfigType = "json"
	Yaml       ConfigType = "yaml"
	Yml        ConfigType = "yml"
)

func (t ConfigType) String() string {
	return string(t)
}

var (
	SupportedConfigTypes = []string{Properties.String(), Json.String(), Yaml.String(), Yml.String()}
)

type apollo struct {
	appConfig            *config.AppConfig
	configType           string
	onConfigChange       func(*storage.ChangeEvent)
	onNewestConfigChange func(*storage.FullChangeEvent)
}

var instance = &apollo{
	appConfig: &config.AppConfig{
		AppID:             "",
		Cluster:           "default",
		NamespaceName:     "application",
		IP:                "",
		IsBackupConfig:    false,
		BackupConfigPath:  "",
		Secret:            "",
		SyncServerTimeout: 0,
		MustStart:         true,
	},
	configType: Properties.String(),
}

// WithAppID apollo 应用名称
func WithAppID(appID string) {
	instance.appConfig.AppID = appID
}

// WithCluster apollo 集群名称
func WithCluster(cluster string) {
	instance.appConfig.Cluster = cluster
}

// WithIsBackupConfig 配置是否备份
func WithIsBackupConfig(isBackupConfig bool) {
	instance.appConfig.IsBackupConfig = isBackupConfig
}

// WithBackupConfigPath 配置备份文件的路径
func WithBackupConfigPath(backupConfigPath string) {
	instance.appConfig.BackupConfigPath = backupConfigPath
}

// WithSyncServerTimeout sync 超时时间
func WithSyncServerTimeout(syncServerTimeout int) {
	instance.appConfig.SyncServerTimeout = syncServerTimeout
}

// WithMustStart 第一次同步是否必须成功
func WithMustStart(mustStart bool) {
	instance.appConfig.MustStart = mustStart
}

// WithConfigType 设置配置文件的格式
func WithConfigType(configType string) {
	instance.configType = configType
}

// WithOnConfigChange 远程配置更新时的回调函数（获取更改的配置）
func WithOnConfigChange(f func(*storage.ChangeEvent)) {
	instance.onConfigChange = f
}

// WithOnNewestConfigChange 远程配置更新时的回调函数（可获取全部配置）
func WithOnNewestConfigChange(f func(*storage.FullChangeEvent)) {
	instance.onNewestConfigChange = f
}

// SetSignature 设置自定义 http 授权控件
func SetSignature(auth auth.HTTPAuth) {
	agollo.SetSignature(auth)
}

// SetBackupFileHandler 设置自定义备份文件处理组件
func SetBackupFileHandler(file file.FileHandler) {
	agollo.SetBackupFileHandler(file)
}

// SetLoadBalance 设置自定义负载均衡组件
func SetLoadBalance(loadBalance cluster.LoadBalance) {
	agollo.SetLoadBalance(loadBalance)
}

// SetLogger 设置自定义logger组件
func SetLogger(loggerInterface log.LoggerInterface) {
	agollo.SetLogger(loggerInterface)
}

// SetCache 设置自定义cache组件
func SetCache(cacheFactory agcache.CacheFactory) {
	agollo.SetCache(cacheFactory)
}

func newClient() (agollo.Client, error) {
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return instance.appConfig, nil
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

type configManager struct {
	client agollo.Client
}

// NewConfigManager 返回 crypt.ConfigManager 的实现
func NewConfigManager(endpoint, namespace, secret string) (crypt.ConfigManager, error) {
	if instance.appConfig.AppID == "" {
		return nil, ErrApolloAppIDMissing
	}

	if !stringInSlice(instance.configType, SupportedConfigTypes) {
		return nil, fmt.Errorf("unsupported config type %s", instance.configType)
	}

	instance.appConfig.IP = endpoint
	instance.appConfig.NamespaceName = namespace
	instance.appConfig.Secret = secret

	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return configManager{client}, nil
}

// Get 获取 properties 格式的配置内容
func (m configManager) Get(key string) ([]byte, error) {
	cache := m.client.GetConfigCache(key)
	content := make(map[string]string, cache.EntryCount())
	cache.Range(func(key, value interface{}) bool {
		content[key.(string)] = fmt.Sprint(value)
		return true
	})

	return marshalToConfigType(content, instance.configType)
}

// Watch 监听配置中心的配置变化
func (m configManager) Watch(key string, stop chan bool) <-chan *crypt.Response {
	resp := make(chan *crypt.Response, 0)
	backendResp := make(chan *listenerResponse, 0)
	m.client.AddChangeListener(listener{backendResp, instance.onConfigChange, instance.onNewestConfigChange})
	go func() {
		for {
			select {
			case <-stop:
				return
			case r := <-backendResp:
				if r.Error != nil {
					resp <- &crypt.Response{
						Value: nil,
						Error: r.Error,
					}
					continue
				}

				resp <- &crypt.Response{
					Value: r.Content,
					Error: nil,
				}
			}
		}
	}()

	return resp
}

func (m configManager) List(key string) (crypt.KVPairs, error) {
	panic("method List is not implemented")
}

func (m configManager) Set(key string, value []byte) error {
	panic("method Set is not implemented")
}

func marshalToConfigType(m map[string]string, configType string) ([]byte, error) {
	var ret []byte

	switch configType {
	case Properties.String():
		p := properties.LoadMap(m)
		ret = []byte(p.String())
	case Json.String():
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		ret = b
	case Yml.String(), Yaml.String():
		b, err := yaml.Marshal(m)
		if err != nil {
			return nil, err
		}
		ret = b
	}

	return ret, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
