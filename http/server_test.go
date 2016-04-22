package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"testing"
	"time"
)

type mockHandler struct {
	called bool
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.called = true
}

func TestSimpleServer(t *testing.T) {
	config := Config{
		Host: "localhost",
		Port: 8080,
	}
	r := mux.NewRouter()
	h := new(mockHandler)
	r.Handle("/", h).Methods("GET")
	s := &SimpleServer{
		Config:  config,
		Handler: r,
	}
	go s.Serve()
	time.Sleep(5 * time.Second)
	http.Get("http://localhost:8080/")
	if !h.called {
		t.Fail()
	}
}
