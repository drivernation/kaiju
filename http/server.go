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

// standardServer is a Server implementation that uses the server from the stdlib's http package to server HTTP requests.
// standardServer can serve either TLS or non-TLS requests depending on configuration.
type standardServer struct {
	Config  Config
	Handler http.Handler
}

func (s *standardServer) Serve() (err error) {
	addr := s.Addr()
	if s.Config.Tls.Enabled {
		err = http.ListenAndServeTLS(addr, s.Config.Tls.Cert, s.Config.Tls.Key, s.Handler)
	} else {
		err = http.ListenAndServe(addr, s.Handler)
	}
	return
}

func (s standardServer) Addr() string {
	return fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
}

// StandardServer creates a new standardServer instance for serving HTTP requests using the stdlib's http package.
func StandardServer(config Config, handler http.Handler) Server {
	return &standardServer{
		Config:  config,
		Handler: handler,
	}
}
