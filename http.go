package kaiju
import (
	"net/http"
	"net/http/httputil"
	"github.com/goods/httpbuf"
)

//Write an empty JSON Response to w. RespCode is provided as the response code.
func EmptyJsonResponse(w http.ResponseWriter, respCode int) {
	w.WriteHeader(respCode)
	w.Header().Add("Content-type", "application/json")
	w.Write([]byte{"{ }"})
}

type loggedHandler func(w http.ResponseWriter, req *http.Request)

func (h loggedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if dump, err := httputil.DumpRequest(req, true); err != nil {
		logger.Errorf("Failed to log request:\n%s", err)
	} else {
		logger.Info(string(dump))
	}

	buf := new(httpbuf.Buffer)
	h(buf, req)
	logger.Info("response: %s", buf.String())
	buf.Apply(w)
}