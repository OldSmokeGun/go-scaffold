package pkg

import (
	"github.com/google/wire"

	"go-scaffold/internal/app/pkg/casbin"
	"go-scaffold/internal/app/pkg/client"
	"go-scaffold/internal/app/pkg/db"
	"go-scaffold/internal/app/pkg/discovery"
	"go-scaffold/internal/app/pkg/ent"
	"go-scaffold/internal/app/pkg/gorm"
	"go-scaffold/internal/app/pkg/redis"
	"go-scaffold/internal/app/pkg/uid"
)

var ProviderSet = wire.NewSet(
	casbin.Provide,
	client.ProvideGRPC,
	db.Provide,
	discovery.Provide,
	ent.Provide,
	gorm.Provide,
	redis.Provide,
	uid.Provide,
)
