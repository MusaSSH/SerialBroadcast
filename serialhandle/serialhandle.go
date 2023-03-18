package serialhandle

import (
	"bytes"
	"context"
	"io"

	"github.com/MusaSSH/SerialBroadcast/config"
	"github.com/MusaSSH/SerialBroadcast/message"
	"go.bug.st/serial"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SerialPort struct {
	port    serial.Port
	buff    bytes.Buffer
	message message.Message
	logger  *zap.Logger
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

			err = s.message.Publish(b)
			if err != nil {
				s.logger.Error("Error publishing message", zap.Error(err))
			}
		}
	}
}

func Build() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, c config.Config, m message.Message, l *zap.Logger) (SerialPort, error) {

		port, err := serial.Open(c.SerialPort, &serial.Mode{
			BaudRate: c.BaudRate,
		})

		if err != nil {
			return SerialPort{}, err
		}
		s := SerialPort{
			logger:  l,
			message: m,
			port:    port,
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
