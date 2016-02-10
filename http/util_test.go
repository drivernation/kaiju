package http

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestEmptyJsonResponse(t *testing.T) {
	w := httptest.NewRecorder()
	respCode := 200
	EmptyJsonResponse(w, respCode)
	assert.Equal(t, respCode, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-type"))
	assert.Equal(t, "{ }", w.Body.String())
}
