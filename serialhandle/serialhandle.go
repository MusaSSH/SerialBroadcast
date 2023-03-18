package serialhandle

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/MusaSSH/SerialBroadcast/config"
	"github.com/MusaSSH/SerialBroadcast/message"
	"go.bug.st/serial"
	"go.uber.org/fx"
)

type SerialPort struct {
	port    serial.Port
	buff    bytes.Buffer
	message message.Message
}

func (s SerialPort) read(c context.Context) {
	for c.Err() == nil {
		read := make([]byte, 128)
		_, err := s.port.Read(read)
		if err != nil {
			log.Println(err)
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
			s.message.Publish(b)
		}
	}
}

func Build() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, c config.Config, m message.Message) (s SerialPort, err error) {
		port, err := serial.Open(c.SerialPort, &serial.Mode{
			BaudRate: c.BaudRate,
		})

		if err != nil {
			return s, err
		}
		s.port = port
		s.message = m

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
		return
	})
}
