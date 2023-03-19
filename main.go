package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/MusaSSH/SerialBroadcast/config"
	"github.com/MusaSSH/SerialBroadcast/logger"
	"github.com/MusaSSH/SerialBroadcast/serialhandle"
	"github.com/MusaSSH/SerialBroadcast/server"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()

	ctx, cf := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cf()

	app := fx.New(
		config.Build(),
		logger.Build(),
		server.Build(),
		serialhandle.Build(),
		fx.Invoke(func(s serialhandle.SerialPort) {}),
	)

	if err := app.Start(ctx); err != nil {
		zap.L().Fatal("Error starting application", zap.Error(err))
	}

	<-ctx.Done()
	ctx, cf = context.WithTimeout(context.Background(), time.Second*10)
	defer cf()

	if err := app.Stop(ctx); err != nil {
		zap.L().Fatal("Error stopping application", zap.Error(err))
	}
}
