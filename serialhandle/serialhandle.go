package serialhandle

import (
	"bytes"
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

func (s SerialPort) read() {
	for {
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

func new(c config.Config) (s SerialPort, err error) {
	port, err := serial.Open(c.SerialPort, &serial.Mode{
		BaudRate: c.BaudRate,
	})

	if err != nil {
		return s, err
	}
	s.port = port
	return
}

func start(s SerialPort) {
	go s.read()
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(new),
		fx.Invoke(start),
	)
}
