package config

import (
	"errors"
	"fmt"

	"github.com/google/wire"
)

var ErrEntryNotConfigured = errors.New("the configuration entry is not configured")

var ProviderSet = wire.NewSet(
	GetApp,
	GetHTTP,
	GetGRPC,
	GetDB,
	GetDBConn,
	GetRedis,
	GetTrace,
	GetDiscovery,
	GetServices,
	GetJWT,
	GetCasbin,
	GetKafka,
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
	DB        *DB        `json:"db"`
	Redis     *Redis     `json:"redis"`
	Trace     *Trace     `json:"trace"`
	Discovery *Discovery `json:"discovery"`
	Services  *Services  `json:"services"`
	JWT       *JWT       `json:"jwt"`
	Casbin    *Casbin    `json:"casbin"`
	Kafka     *Kafka     `json:"kafka"`
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

func GetHTTP() (HTTP, error) {
	return getEntry(config.HTTP)
}

func GetGRPC() (GRPC, error) {
	return getEntry(config.GRPC)
}

func GetDB() (DB, error) {
	return getEntry(config.DB)
}

func GetDBConn() (DBConn, error) {
	dbConfig, err := GetDB()
	if err != nil {
		return DBConn{}, err
	}
	return dbConfig.DBConn, nil
}

func GetRedis() (Redis, error) {
	return getEntry(config.Redis)
}

func GetTrace() (Trace, error) {
	return getEntry(config.Trace)
}

func GetDiscovery() (Discovery, error) {
	return getEntry(config.Discovery)
}

func GetServices() (Services, error) {
	return getEntry(config.Services)
}

func GetJWT() (JWT, error) {
	return getEntry(config.JWT)
}

func GetCasbin() (Casbin, error) {
	return getEntry(config.Casbin)
}

func GetKafka() (Kafka, error) {
	return getEntry(config.Kafka)
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
