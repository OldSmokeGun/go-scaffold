package scripts

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	// scripts
	NewExampleCmd,
)
