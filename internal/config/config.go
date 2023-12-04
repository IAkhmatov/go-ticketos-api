//go:generate go-enum --marshal --mustparse

package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

// ENUM(local, tests, production).
type AppEnv string

type Config struct {
	appEnv AppEnv `default:"local" env:"APP_ENV"`

	// DBConnectString PSQL connection string
	DBConnectString string `default:"postgresql://postgres:postgres@localhost:5435/ticketos" env:"DB_CONNECT_STRING"`
	// TestDBConnectString PSQL connection string.
	// uses only in tests
	TestDBConnectString string `default:"postgresql://postgres:postgres@localhost:5435/test" env:"TEST_DB_CONNECT_STRING"`

	// APIReadTimeout API read timeout
	APIReadTimeout int `default:"1"`
	// ApiTimeout API timeout for graceful shutdown
	APITimeout int `default:"5" env:"API_TIMEOUT"`
	// APIPort API port
	APIPort int `default:"7565" env:"API_PORT"`

	// PayMasterTimeout timeout for requests to paymaster
	PayMasterTimeout int `default:"5" env:"PAY_MASTER_TIMEOUT"`
	// PayMasterAPIKey API key for paymaster
	PayMasterAPIKey string `default:"test" env:"PAY_MASTER_API_KEY"`
	// PayMasterMerchantID merchant id for paymaster
	PayMasterMerchantID string `default:"54a1a08f-b344-44b1-9f06-fe5868170e5f" env:"PAY_MASTER_MERCHANT_ID"`
	// PayMasterCallbackURL callback url for paymaster.
	// On this url paymaster will send webhook about invoice
	PayMasterCallbackURL string `default:"https://test.test" env:"PAY_MASTER_CALLBACK_URL"`
	// PayMasterReturnURL return url for paymaster
	// On this url paymaster will redirect user after payment
	PayMasterReturnURL string `default:"https://test.test" env:"PAY_MASTER_RETURN_URL"`

	// OrderTTL ttl for order
	// After this time order will be expired
	OrderTTL int `default:"5" env:"ORDER_TTL"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := configor.Load(&cfg); err != nil {
		return nil, fmt.Errorf("NewConfig: can not load config: %w", err)
	}
	return &cfg, nil
}

func (c Config) IsProduction() bool {
	return c.appEnv == AppEnvProduction
}

func (c Config) IsTests() bool {
	return c.appEnv == AppEnvTests
}
