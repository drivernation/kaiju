package kaiju

import (
	"github.com/goods/httpbuf"
	"net/http"
	"net/http/httputil"
)

//Write an empty JSON Response to w. RespCode is provided as the response code.
func EmptyJsonResponse(w http.ResponseWriter, respCode int) {
	w.WriteHeader(respCode)
	w.Header().Add("Content-type", "application/json")
	w.Write([]byte("{ }"))
}

// A loggedHandler is an http.Handler implementation that logs the request before using the current logger before
// servicing the request.
//
// After servicing the request, the response is logged using the same logger, before the response is ultimately sent
// to the downstream client.
type loggedHandler func(w http.ResponseWriter, req *http.Request)

func (h loggedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if dump, err := httputil.DumpRequest(req, true); err != nil {
		logger.Errorf("Failed to log request:\n%s", err)
	} else {
		logger.Info(string(dump))
	}

	buf := new(httpbuf.Buffer)
	h(buf, req)
	logger.Infof("response: %s", buf.String())
	buf.Apply(w)
}
