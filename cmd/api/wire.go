//go:build wireinject

package main

import (
	"errors"

	"go-ticketos/internal/config"
	"go-ticketos/internal/service"
	"go-ticketos/internal/storage"
	"go-ticketos/internal/transport"
	"go-ticketos/pkg/log"
	"go-ticketos/pkg/paymasterclient"

	"github.com/google/wire"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type API struct {
	eventController     transport.EventController[fiber.Ctx]
	orderController     transport.OrderController[fiber.Ctx]
	paymasterController transport.PayMasterController[fiber.Ctx]
	logger              *zerolog.Logger
	cfg                 *config.Config
}

func NewAPI(
	eventController transport.EventController[fiber.Ctx],
	orderController transport.OrderController[fiber.Ctx],
	paymasterController transport.PayMasterController[fiber.Ctx],
	logger *zerolog.Logger,
	cfg *config.Config,
) (*API, error) {
	if eventController == nil {
		return nil, errors.New("NewAPI: eventController is nil")
	}
	if orderController == nil {
		return nil, errors.New("NewAPI: orderController is nil")
	}
	if paymasterController == nil {
		return nil, errors.New("NewAPI: paymasterController is nil")
	}
	if logger == nil {
		return nil, errors.New("NewAPI: logger is nil")
	}
	return &API{
		eventController:     eventController,
		orderController:     orderController,
		paymasterController: paymasterController,
		logger:              logger,
		cfg:                 cfg,
	}, nil
}

func InitAPI() (*API, error) {
	wire.Build(
		config.ConfigProviderSet,
		log.LoggerProviderSet,
		storage.SqlxDBProviderSet,
		storage.EventRepoProviderSet,
		service.EventServiceProviderSet,
		transport.EventControllerProviderSet,
		storage.TicketCategoryRepoProviderSet,
		storage.PromocodeRepoProviderSet,
		storage.OrderRepoProviderSet,
		paymasterclient.PayMasterClientProviderSet,
		service.OrderServiceProviderSet,
		transport.OrderControllerProviderSet,
		service.PayMasterWebHookUseCaseProviderSet,
		transport.PayMasterControllerProviderSet,
		NewAPI,
	)
	return &API{}, nil
}
