package paymasterclient

import (
	"github.com/google/wire"
)

var PayMasterClientProviderSet wire.ProviderSet = wire.NewSet(
	NewPayMasterClient,
	wire.Bind(new(PayMasterClient), new(*payMasterClient)),
)
