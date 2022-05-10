//go:build wireinject
// +build wireinject

package tests

import "github.com/google/wire"

func Init() (*Tests, func(), error) {
	panic(wire.Build(
		ProviderSet,
		New,
	))
}
