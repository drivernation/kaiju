package http

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggedHandler(t *testing.T) {
	hit := false
	h := LoggedHandler{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			hit = true
			w.Write([]byte("blah"))
		}),
	}
	server := httptest.NewServer(h)
	defer server.Close()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)
	assert.True(t, hit)
	assert.Equal(t, "blah", string(body))
}
