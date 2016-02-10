package http

import (
	"net/http"
)

//Write an empty JSON Response to w. RespCode is provided as the response code.
func EmptyJsonResponse(w http.ResponseWriter, respCode int) {
	w.WriteHeader(respCode)
	w.Header().Add("Content-type", "application/json")
	w.Write([]byte("{ }"))
}
