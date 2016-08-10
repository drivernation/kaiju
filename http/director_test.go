package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoundRobin_Direct(t *testing.T) {
	rr := RoundRobin([]string{"test1", "test2"}, SchemeHttp)
	req, err := rr.Direct("get", "/", []byte("blah"))
	assert.NoError(t, err)
	assert.Contains(t, req.URL.String(), "test1")
	req, err = rr.Direct("get", "/", nil)
	assert.NoError(t, err)
	assert.Contains(t, req.URL.String(), "test2")
}
