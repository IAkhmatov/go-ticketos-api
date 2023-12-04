package config

import "github.com/google/wire"

var ConfigProviderSet wire.ProviderSet = wire.NewSet(
	NewConfig,
)
