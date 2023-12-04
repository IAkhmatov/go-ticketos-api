package service

import (
	"go-ticketos/internal/domain"

	"github.com/google/wire"
)

var EventServiceProviderSet wire.ProviderSet = wire.NewSet(
	NewEventService,
	wire.Bind(new(domain.EventService), new(*eventService)),
)

var OrderServiceProviderSet wire.ProviderSet = wire.NewSet(
	NewOrderService,
	wire.Bind(new(domain.OrderService), new(*orderService)),
)
var PayMasterWebHookUseCaseProviderSet wire.ProviderSet = wire.NewSet(
	NewPayMasterWebHookUseCase,
	wire.Bind(new(domain.PayMasterWebHookUseCase), new(*payMasterWebHookUseCase)),
)
