package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Production bool          `envconfig:"PRODUCTION" default:"false"`
	LogLevel   zapcore.Level `envconfig:"LOG_LEVEL" default:"info"`

	SPort     string `envconfig:"SERIAL_PORT" required:"true"`
	SBaudRate int    `envconfig:"SERIAL_BAUD_RATE" default:"9600"`

	WSPort string `envconfig:"WEBSOCKET_PORT" default:":1234"`
}

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, err
		}
		return
	})
}
