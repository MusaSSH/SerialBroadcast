package server

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/MusaSSH/SerialBroadcast/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

type Server struct {
	*sync.RWMutex
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

	s.Lock()
	s.conns[c] = true
	s.Unlock()

	for {
		_, _, err := c.Read(r.Context())
		if err != nil {
			s.Lock()
			delete(s.conns, c)
			c.Close(websocket.StatusAbnormalClosure, "")
			s.Unlock()
			return
		}
	}
}

func (s Server) Broadcast(data []byte) error {
	for c := range s.conns {
		err := c.Write(context.Background(), websocket.MessageText, data)
		if err != nil {
			s.Lock()
			delete(s.conns, c)
			c.Close(websocket.StatusInternalError, "")
			s.Unlock()
			return err
		}
	}
	return nil
}

func NewHTTPServer(lc fx.Lifecycle, l *zap.Logger, c config.Config) *http.Server {
	srv := &http.Server{
		Handler: &Server{
			RWMutex: &sync.RWMutex{},
			logger:  l,
			conns:   make(map[*websocket.Conn]bool),
		},
	}

	sctx, cf := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			ln, err := net.Listen("tcp", c.WSPort)
			if err != nil {
				return err
			}

			srv.BaseContext = func(_ net.Listener) context.Context { return sctx }

			go func() {
				if err := srv.Serve(ln); err != nil {
					l.Fatal("Error serving websocket", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			cf()
			return nil
		},
	})
	return srv
}

func Build() fx.Option {
	return fx.Provide(NewHTTPServer)
}
