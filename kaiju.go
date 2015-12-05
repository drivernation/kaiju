package kaiju

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"os/signal"
)

var serviceManager *ServiceManager
var muxer *mux.Router
var tlsListener net.Listener

func init() {
	serviceManager = NewServiceManager()
	muxer = mux.NewRouter()
	registerControlCFunction()
}

// Manages one or more Services. The lifecycles of managed Services are controlled by Kaiju, so the application
// should never start or stop a service manually.
func Manage(services ...Service) {
	for _, s := range services {
		serviceManager.AddService(s)
	}
}

// Registers a http.Handler with Kaiju's request muxer.
func Handle(path string, handler http.Handler, methods ...string) {
	HandleFunc(path, handler.ServeHTTP, methods...)
}

// Registers a handler function with Kaiju's HTTP muxer.
func HandleFunc(path string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	logger.Infof("Handling method(s) %s at %s", methods, path)
	muxer.Handle(path, loggedHandler(handler)).Methods(methods...)
}

// Starts Kaiju using a standard HTTP server. In this mode, HTTP requests are unencrypted.
// An error is returned it kaiju fails to start. This could be due to the ServiceManager encountering an error or if
// the HTTP server failed to start up.
func Start(config Config) error {
	logger.Infof("Starting %d services...", serviceManager.Size())
	if err := serviceManager.Start(); err != nil {
		return err
	}
	defer stopServices()
	http.Handle("/", muxer)
	addr := fmt.Sprintf("%s:%d", config.BindHost, config.Port)
	logger.Infof("Listening on %s.", addr)
	return http.ListenAndServe(addr, nil)
}

// Starts Kaiju using a HTTPS server. In this mode, HTTP requests are encrypted using the provided TLS configuration.
// An error is returned it kaiju fails to start. This could be due to the ServiceManager encountering an error or if
// the HTTPS server failed to start up.
func StartTLS(config Config, tlsConfig *tls.Config) error {
	logger.Infof("Starting %d services...", serviceManager.Size())
	if err := serviceManager.Start(); err != nil {
		return err
	}
	defer stopServices()
	http.Handle("/", muxer)
	addr := fmt.Sprintf("%s:%d", config.BindHost, config.Port)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tlsListener = tls.NewListener(conn, tlsConfig)
	logger.Infof("Listening on %s.", addr)
	return http.Serve(tlsListener, nil)
}

func stopServices() {
	logger.Infof("Stopping %d services...", serviceManager.Size())
	if err := serviceManager.Stop(); err != nil {
		logger.Errorf("Failed to stop one or more services: %s", err)
	}
}

func registerControlCFunction() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			stopServices()
			logger.Close()
			os.Exit(1)
		}
	}()
}
