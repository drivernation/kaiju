//+build !test

package kaiju

import (
	"github.com/Sirupsen/logrus"
	"github.com/drivernation/kaiju/http"
	"github.com/drivernation/kaiju/logging"
	"github.com/gorilla/mux"
	"os"
	"os/signal"
)

var (
	Muxer         *mux.Router = mux.NewRouter()
	shutdownHooks []func()    = make([]func(), 0)
)

func init() {
	RegisterControlCFunction()
}

type Config struct {
	Http    http.Config
	Logging logging.Config
}

// AddShutdownHook appends a function to the list of functions that will be executed before the application shuts down.
// Note that these hooks will not be executed in the event of a panic.
func AddShutdownHook(f func()) {
	shutdownHooks = append(shutdownHooks, f)
}

// Starts Kaiju using a kaiju/http.SimpleHttpServer with the provided config.
// An error is returned it kaiju fails to start.
func Start(config Config) error {
	closer := logging.Configure(config.Logging)
	AddShutdownHook(func() {
		if err := closer(); err != nil {
			logrus.Error(err)
		}
	})
	// Execute any registered shutdown hooks.
	defer func() {
		for _, f := range shutdownHooks {
			f()
		}
	}()
	server := http.StandardServer(config.Http, &http.LoggedHandler{Muxer})
	addr := server.Addr()
	logrus.Infof("Listening on %s.", addr)
	return server.Serve()
}

// Registers a signal handler to kill the app if Ctrl-C is pressed.
func RegisterControlCFunction() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			for _, f := range shutdownHooks {
				f()
			}
			os.Exit(0)
		}
	}()
}
