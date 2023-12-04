package log

import (
	"github.com/google/wire"
)

var LoggerProviderSet wire.ProviderSet = wire.NewSet(
	NewLogger,
)
