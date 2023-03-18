package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/MusaSSH/SerialBroadcast/config"
	"github.com/MusaSSH/SerialBroadcast/message"
	"github.com/MusaSSH/SerialBroadcast/serialhandle"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load()

	ctx, cf := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cf()

	app := fx.New(
		config.Build(),
		message.Build(),
		serialhandle.Build(),
		fx.Invoke(func(s serialhandle.SerialPort) {}),
	)
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	ctx, cf = context.WithTimeout(context.Background(), time.Second*10)
	defer cf()

	if err := app.Stop(ctx); err != nil {
		log.Fatal(err)
	}
}
