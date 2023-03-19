package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

type Config struct {
	SerialPort string `envconfig:"SERIAL_PORT" required:"true"`
	BaudRate   int    `envconfig:"BAUD_RATE" default:"9600"`
	WSPort     string `envconfig:"WS_PORT" default:":1234"`
}

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, err
		}
		return
	})
}
