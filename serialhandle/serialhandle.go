package serialhandle

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/MusaSSH/SerialBroadcast/config"
	"go.bug.st/serial"
	"go.uber.org/fx"
)

type SerialPort struct {
	port serial.Port
	buff bytes.Buffer
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
			fmt.Print(string(b))
		}
	}
}

func Build() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, c config.Config) (s SerialPort, err error) {
		port, err := serial.Open(c.SerialPort, &serial.Mode{
			BaudRate: c.BaudRate,
		})

		if err != nil {
			return s, err
		}
		s.port = port

		stopc, stop := context.WithCancel(context.Background())

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go s.read(stopc)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				stop()
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
