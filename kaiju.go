//+build !test

package kaiju

import (
	"github.com/drivernation/kaiju/http"
	"github.com/drivernation/kaiju/logging"
	"github.com/gorilla/mux"
	"os"
	"os/signal"
)

var (
	Muxer        *mux.Router = mux.NewRouter()
	ShutdownHook func()
)

func init() {
	RegisterControlCFunction()
}

type Config struct {
	Http http.Config
}

// Starts Kaiju using a kaiju/http.SimpleHttpServer with the provided config.
// An error is returned it kaiju fails to start.
func Start(config Config) error {
	server := &http.SimpleServer{
		Config:  config.Http,
		Handler: Muxer,
	}
	addr := server.Addr()
	logging.Logger.Infof("Listening on %s.", addr)
	return server.Serve()
}

// Registers a signal handler to kill the app if Ctrl-C is pressed.
func RegisterControlCFunction() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logging.Logger.Close()
			if ShutdownHook != nil {
				ShutdownHook()
			}
			os.Exit(0)
		}
	}()
}
