package remote

import (
	"bytes"
	crypt "github.com/sagikazarmark/crypt/config"
	"github.com/spf13/viper"
	"go-scaffold/pkg/config/remote/apollo"
	"io"
	"os"
)

func init() {
	viper.SupportedRemoteProviders = append(
		viper.SupportedRemoteProviders,
		"apollo",
	)
	viper.RemoteConfig = &remoteConfigProvider{}
}

type remoteConfigProvider struct{}

func (rc remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	b, err := cm.Get(rp.Path())
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (rc remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	resp, err := cm.Get(rp.Path())
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(resp), nil
}

func (rc remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, nil
	}
	quit := make(chan bool)
	quitwc := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)
	cryptoResponseCh := cm.Watch(rp.Path(), quit)
	// need this function to convert the Channel response form crypt.Response to viper.Response
	go func(cr <-chan *crypt.Response, vr chan<- *viper.RemoteResponse, quitwc <-chan bool, quit chan<- bool) {
		for {
			select {
			case <-quitwc:
				quit <- true
				return
			case resp := <-cr:
				vr <- &viper.RemoteResponse{
					Error: resp.Error,
					Value: resp.Value,
				}
			}
		}
	}(cryptoResponseCh, viperResponsCh, quitwc, quit)

	return viperResponsCh, quitwc
}

func getConfigManager(rp viper.RemoteProvider) (crypt.ConfigManager, error) {
	var cm crypt.ConfigManager
	var err error

	if rp.SecretKeyring() != "" {
		// apollo 直接使用 secret
		switch rp.Provider() {
		case "apollo":
			cm, err = apollo.NewConfigManager(rp.Endpoint(), rp.Path(), rp.SecretKeyring())
			if err != nil {
				return nil, err
			}
			return cm, nil
		}

		var kr *os.File
		kr, err = os.Open(rp.SecretKeyring())
		if err != nil {
			return nil, err
		}
		defer kr.Close()
		switch rp.Provider() {
		case "etcd":
			cm, err = crypt.NewEtcdConfigManager([]string{rp.Endpoint()}, kr)
		case "firestore":
			cm, err = crypt.NewFirestoreConfigManager([]string{rp.Endpoint()}, kr)
		default:
			cm, err = crypt.NewConsulConfigManager([]string{rp.Endpoint()}, kr)
		}
	} else {
		switch rp.Provider() {
		case "etcd":
			cm, err = crypt.NewStandardEtcdConfigManager([]string{rp.Endpoint()})
		case "firestore":
			cm, err = crypt.NewStandardFirestoreConfigManager([]string{rp.Endpoint()})
		case "apollo":
			cm, err = apollo.NewConfigManager(rp.Endpoint(), rp.Path(), "")
		default:
			cm, err = crypt.NewStandardConsulConfigManager([]string{rp.Endpoint()})
		}
	}
	if err != nil {
		return nil, err
	}
	return cm, nil
}
