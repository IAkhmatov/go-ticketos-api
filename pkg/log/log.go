package log

import (
	"os"
	"time"

	"go-ticketos/internal/config"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *config.Config) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}

	logger := zerolog.New(output).With().Timestamp().Caller().Stack().Logger()

	if cfg.IsTests() {
		logger = logger.Level(zerolog.Disabled)
	}
	return &logger
}
