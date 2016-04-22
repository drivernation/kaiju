package http

import (
	"fmt"
	"net/http"
)

type Server interface {
	Serve() error
	Addr() string
}

type Config struct {
	Host string
	Port int
	Tls  struct {
		Enabled bool
		Cert    string
		Key     string
	}
}

type SimpleServer struct {
	Config  Config
	Handler http.Handler
}

func (s *SimpleServer) Serve() (err error) {
	addr := s.Addr()
	if s.Config.Tls.Enabled {
		err = http.ListenAndServeTLS(addr, s.Config.Tls.Cert, s.Config.Tls.Key, s.Handler)
	} else {
		err = http.ListenAndServe(addr, s.Handler)
	}
	return
}

func (s SimpleServer) Addr() string {
	return fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
}
