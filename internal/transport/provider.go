package transport

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var EventControllerProviderSet wire.ProviderSet = wire.NewSet(
	NewEventController,
	wire.Bind(new(EventController[fiber.Ctx]), new(*eventControllerFiber)),
)

var OrderControllerProviderSet wire.ProviderSet = wire.NewSet(
	NewOrderController,
	wire.Bind(new(OrderController[fiber.Ctx]), new(*orderControllerFiber)),
)

var PayMasterControllerProviderSet wire.ProviderSet = wire.NewSet(
	NewPayMasterController,
	wire.Bind(new(PayMasterController[fiber.Ctx]), new(*payMasterControllerFiber)),
)
