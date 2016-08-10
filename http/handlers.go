package http

import (
	"github.com/Sirupsen/logrus"
	"github.com/goods/httpbuf"
	"net/http"
	"net/http/httputil"
)

// A loggedHandler is an http.Handler implementation that logs the request before using the current logger before
// servicing the request.
//
// After servicing the request, the response is logged using the same logger, before the response is ultimately sent
// to the downstream client.
type LoggedHandler func(w http.ResponseWriter, req *http.Request)

func (h LoggedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if dump, err := httputil.DumpRequest(req, true); err != nil {
		logrus.Errorf("Failed to log request:\n%s", err)
	} else {
		logrus.Info(string(dump))
	}

	buf := new(httpbuf.Buffer)
	h(buf, req)
	logrus.Infof("response: %s", buf.String())
	buf.Apply(w)
}
