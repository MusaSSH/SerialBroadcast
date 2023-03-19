package server

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

type Server struct {
	conns  map[*websocket.Conn]bool
	logger *zap.Logger
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		s.logger.Error("Error accepting websocket", zap.Error(err))
		return
	}

	s.conns[c] = true
}

func (s Server) Broadcast(data []byte) error {
	for c := range s.conns {
		err := c.Write(context.Background(), websocket.MessageText, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewHTTPServer(l *zap.Logger) *http.Server {
	return &http.Server{
		Handler: &Server{
			logger: l,
			conns:  make(map[*websocket.Conn]bool),
		},
	}
}

func StartHTTPServer(s *http.Server) {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go func() {
		s.Serve(l)
	}()

}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(NewHTTPServer),
		fx.Invoke(StartHTTPServer),
	)
}
