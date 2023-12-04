// nolint: cyclop
package main

import (
	"time"

	"go-ticketos/internal/transport"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
)

func main() {
	api, err := InitAPI()
	if err != nil {
		panic(err)
	}
	api.logger.Info().Msg("Building app")
	app := fiber.New(
		fiber.Config{
			ReadTimeout: time.Duration(api.cfg.APIReadTimeout) * time.Second,
		},
	)
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: api.logger,
	}))
	apiGroup := app.Group("/api")
	v1 := apiGroup.Group("/v1")
	v1.Get("live", transport.Live)
	v1.Get("ready", transport.Ready)

	event := v1.Group("/event")
	event.Post("", api.eventController.Create)

	order := v1.Group("/order")
	order.Post("", api.orderController.Create)

	webhook := v1.Group("/webhook")
	webhook.Post("paymaster", api.paymasterController.WebHook)

	transport.RunAPI(app, api.cfg.APITimeout, api.cfg.APIPort, api.logger)
}
