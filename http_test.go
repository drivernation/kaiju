package kaiju
import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"io/ioutil"
)

func TestEmptyJsonResponse(t *testing.T) {
	w := httptest.NewRecorder()
	respCode := 200
	EmptyJsonResponse(w, respCode)
	assert.Equal(t, respCode, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-type"))
	assert.Equal(t, "{ }", w.Body.String())
}

func TestLoggedHandler(t *testing.T) {
	hit := false
	h := loggedHandler(func(w http.ResponseWriter, req *http.Request) {
		hit=true
		w.Write([]byte("blah"))
	})
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
