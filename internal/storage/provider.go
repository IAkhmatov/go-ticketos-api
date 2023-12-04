package storage

import (
	"go-ticketos/internal/domain"

	"github.com/google/wire"
)

var (
	SqlxDBProviderSet wire.ProviderSet = wire.NewSet(
		NewSqlxDBFromConfig,
	)

	EventRepoProviderSet wire.ProviderSet = wire.NewSet(
		NewEventRepo,
		wire.Bind(new(domain.EventRepo), new(*eventPSQLRepo)),
	)

	OrderRepoProviderSet wire.ProviderSet = wire.NewSet(
		NewOrderRepo,
		wire.Bind(new(domain.OrderRepo), new(*orderPSQLRepo)),
	)

	TicketCategoryRepoProviderSet wire.ProviderSet = wire.NewSet(
		NewTicketCategoryRepo,
		wire.Bind(new(domain.TicketCategoryRepo), new(*ticketCategoryPSQLRepo)),
	)

	PromocodeRepoProviderSet wire.ProviderSet = wire.NewSet(
		NewPromocodeRepo,
		wire.Bind(new(domain.PromocodeRepo), new(*promocodePSQLRepo)),
	)
)
