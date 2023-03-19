package serialhandle

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/MusaSSH/SerialBroadcast/config"
	"github.com/MusaSSH/SerialBroadcast/server"
	"go.bug.st/serial"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SerialPort struct {
	port   serial.Port
	buff   bytes.Buffer
	ws     *http.Server
	logger *zap.Logger
}

func (s SerialPort) read(c context.Context) {
	for c.Err() == nil {
		read := make([]byte, 128)
		_, err := s.port.Read(read)
		if err != nil {
			s.logger.Error("Error reading from serial port", zap.Error(err))
		}
		read = bytes.Trim(read, "\x00")
		s.buff.Write(read)

		for {
			b, err := s.buff.ReadBytes('\n')
			if err == io.EOF {
				s.buff.Reset()
				s.buff.Write(b)
				break
			}
			err = s.ws.Handler.(*server.Server).Broadcast(b)
			if err != nil {
				s.logger.Error("Error broadcasting to websocket", zap.Error(err))
			}
		}
	}
}

func Build() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, c config.Config, l *zap.Logger, ws *http.Server) (SerialPort, error) {

		port, err := serial.Open(c.SerialPort, &serial.Mode{
			BaudRate: c.BaudRate,
		})

		if err != nil {
			return SerialPort{}, err
		}
		s := SerialPort{
			logger: l,
			port:   port,
			ws:     ws,
		}

		sctx, sf := context.WithCancel(context.Background())

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go s.read(sctx)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				sf()
				err := s.port.Close()
				if err != nil {
					return err
				}
				return nil
			},
		})
		return s, nil
	})
}
