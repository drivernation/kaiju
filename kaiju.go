//+build !test

package kaiju

import (
	"crypto/tls"
	"fmt"
	ghttp "github.com/drivernation/kaiju/http"
	"github.com/drivernation/kaiju/logging"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"os/signal"
)

var muxer *mux.Router = mux.NewRouter()
var tlsListener net.Listener

func init() {
	registerControlCFunction()
}

// Registers a http.Handler with Kaiju's request muxer.
func Handle(path string, handler http.Handler, methods ...string) {
	HandleFunc(path, handler.ServeHTTP, methods...)
}

// Registers a handler function with Kaiju's HTTP muxer.
func HandleFunc(path string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	logging.Logger.Infof("Handling method(s) %s at %s", methods, path)
	muxer.Handle(path, ghttp.LoggedHandler(handler)).Methods(methods...)
}

// Starts Kaiju using a standard HTTP server. In this mode, HTTP requests are unencrypted.
// An error is returned it kaiju fails to start.
func Start(config Config) error {
	http.Handle("/", muxer)
	addr := fmt.Sprintf("%s:%d", config.BindHost, config.Port)
	logging.Logger.Infof("Listening on %s.", addr)
	return http.ListenAndServe(addr, nil)
}

// Starts Kaiju using a HTTPS server. In this mode, HTTP requests are encrypted using the provided TLS configuration.
// An error is returned it kaiju fails to start.
func StartTLS(config Config, tlsConfig *tls.Config) error {
	http.Handle("/", muxer)
	addr := fmt.Sprintf("%s:%d", config.BindHost, config.Port)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tlsListener = tls.NewListener(conn, tlsConfig)
	logging.Logger.Infof("Listening on %s.", addr)
	return http.Serve(tlsListener, nil)
}

func registerControlCFunction() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logging.Logger.Close()
			os.Exit(1)
		}
	}()
}
