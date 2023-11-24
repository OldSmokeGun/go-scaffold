package config

import (
	"errors"
	"fmt"

	"github.com/google/wire"
)

var ErrEntryNotConfigured = errors.New("the configuration entry is not configured")

var ProviderSet = wire.NewSet(
	GetApp,
	GetHTTPServer,
	GetHTTPCasbin,
	GetGRPCServer,
	GetServices,
	GetDiscovery,
	GetDatabase,
	GetDatabaseConn,
	GetRedis,
	GetKafka,
	GetTrace,
)

type Configure interface {
	GetName() string
}

var config *Config

// Config application config
type Config struct {
	App       *App       `json:"app"`
	HTTP      *HTTP      `json:"http"`
	GRPC      *GRPC      `json:"grpc"`
	Services  *Services  `json:"services"`
	Discovery *Discovery `json:"discovery"`
	Database  *Database  `json:"database"`
	Redis     *Redis     `json:"redis"`
	Kafka     *Kafka     `json:"kafka"`
	Trace     *Trace     `json:"trace"`
}

// SetConfig set configuration
func SetConfig(c *Config) {
	config = c
}

// HasConfigured return whether the configuration is set
func HasConfigured() bool {
	return config != nil
}

// IsNotConfigured return whether error is ErrEntryNotConfigured
func IsNotConfigured(err error) bool {
	return errors.Is(err, ErrEntryNotConfigured)
}

func GetApp() (App, error) {
	return getEntry(config.App)
}

func GetHTTPServer() (HTTPServer, error) {
	httpConfig, err := getHTTP()
	if err != nil {
		return HTTPServer{}, err
	}
	return getEntry(httpConfig.Server)
}

func GetHTTPCasbin() (Casbin, error) {
	httpConfig, err := getHTTP()
	if err != nil {
		return Casbin{}, err
	}
	return getEntry(httpConfig.Casbin)
}

func GetGRPCServer() (GRPCServer, error) {
	grpcConfig, err := getGRPC()
	if err != nil {
		return GRPCServer{}, err
	}
	return getEntry(grpcConfig.Server)
}

func GetServices() (Services, error) {
	return getEntry(config.Services)
}

func GetDiscovery() (Discovery, error) {
	return getEntry(config.Discovery)
}

func GetDatabase() (Database, error) {
	return getEntry(config.Database)
}

func GetDatabaseConn() (DatabaseConn, error) {
	dbConfig, err := GetDatabase()
	if err != nil {
		return DatabaseConn{}, err
	}
	return dbConfig.DatabaseConn, nil
}

func GetRedis() (Redis, error) {
	return getEntry(config.Redis)
}

func GetKafka() (Kafka, error) {
	return getEntry(config.Kafka)
}

func GetTrace() (Trace, error) {
	return getEntry(config.Trace)
}

func getHTTP() (HTTP, error) {
	return getEntry(config.HTTP)
}

func getGRPC() (GRPC, error) {
	return getEntry(config.GRPC)
}

func getEntry[T Configure](t *T) (T, error) {
	if t == nil {
		e := new(T)
		return *e, wrapEntryNotConfiguredError(*e)
	}
	return *t, nil
}

func wrapEntryNotConfiguredError(c Configure) error {
	return fmt.Errorf("%w: %s", ErrEntryNotConfigured, c.GetName())
}
