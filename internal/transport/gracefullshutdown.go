package transport

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func RunAPI(app *fiber.App, timeout int, port int, logger *zerolog.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	serverShutdown := make(chan struct{}, 1)
	go func() {
		<-c
		logger.Info().Int("timeout", timeout).
			Msg("Gracefully shutdown app. App will be forcefully close all cons after timeout")
		if err := app.ShutdownWithTimeout(time.Duration(timeout) * time.Second); err != nil {
			logger.Panic().Err(err).Msg("Can not finish all requests. Force close app.")
			panic(err)
		}
		serverShutdown <- struct{}{}
	}()

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		logger.Panic().Err(err).Msg("Can not start app")
	}

	<-serverShutdown
	logger.Info().Msg("Close app")
}
