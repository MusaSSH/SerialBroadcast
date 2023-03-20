package logger

import (
	"github.com/MusaSSH/SerialBroadcast/config"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
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
	var zapconfig zap.Config
	if c.Production {
		zapconfig = zap.NewProductionConfig()
	} else {
		zapconfig = zap.NewDevelopmentConfig()
	}

	zapconfig.Level.SetLevel(c.LogLevel)
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapconfig.Build()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to build zap config"))
	}
	return logger, err
}
