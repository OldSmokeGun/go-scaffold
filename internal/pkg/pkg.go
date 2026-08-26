package pkg

import (
	"github.com/google/wire"

	"go-scaffold/internal/pkg/casbin"
	"go-scaffold/internal/pkg/client"
	"go-scaffold/internal/pkg/db"
	"go-scaffold/internal/pkg/discovery"
	"go-scaffold/internal/pkg/gorm"
	"go-scaffold/internal/pkg/redis"
	"go-scaffold/internal/pkg/uid"
)

var ProviderSet = wire.NewSet(
	casbin.Provide,
	client.ProvideGRPC,
	db.Provide,
	discovery.Provide,
	gorm.ProvideDefault,
	redis.ProvideDefault,
	uid.Provide,
)
