package logger

import (
	"github.com/MusaSSH/SerialBroadcast/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newLogger),
	)
}

func newLogger(c config.Config) (*zap.Logger, error) {
	zapconfig := zap.NewDevelopmentConfig()
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapconfig.Build()
	return logger, err
}
