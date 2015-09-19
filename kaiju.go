package kaiju
import (
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"fmt"
	"crypto/tls"
	"os"
	"os/signal"
)

var services []Service
var muxer *mux.Router
var tlsListener net.Listener

func init() {
	services = []Service{}
	muxer = mux.NewRouter()
	registerControlCFunction()
}

func Manage(ss ...Service) {
	services = append(services, ss...)
}

func Handle(path string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	logger.Infof("Handling method(s) %s at %s", methods, path)
	muxer.Handle(path, loggedHandler(handler)).Methods(methods...)
}

func Start(config Config) error {
	logger.Infof("Starting services...")
	err := startServices()
	if err != nil {
		return err
	}
	defer stopServices()
	http.Handle("/", muxer)
	addr := fmt.Sprintf("%s:%d", config.BindHost, config.Port)
	logger.Infof("Listening on %s.", addr)
	return http.ListenAndServe(addr, nil)
}

func StartTLS(config Config, tlsConfig *tls.Config) error {
	logger.Infof("Starting services...")
	if err := startServices(); err != nil {
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

func startServices() error {
	for _, service := range services {
		if err := service.Start(); err != nil {
			return err
		}
	}
	return nil
}

func stopServices() {
	for _, service := range services {
		service.Stop()
	}
}

func registerControlCFunction() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			stopServices()
		}
	}()
}
