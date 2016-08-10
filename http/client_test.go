package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type mockClient struct {
	mock.Mock
}

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	args := c.Mock.Called(req)
	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*http.Response), args.Error(1)
}

type mockDirector struct {
	mock.Mock
}

func (d *mockDirector) Direct(method, resource string, body []byte) (*http.Request, error) {
	args := d.Mock.Called(method, resource, body)
	req := args.Get(0)
	if req == nil {
		return nil, args.Error(1)
	}

	return req.(*http.Request), args.Error(1)
}

func TestDirectedClient_IssueRequest(t *testing.T) {
	mc := new(mockClient)
	req := &http.Request{}
	resp := &http.Response{StatusCode: 200}
	mc.On("Do", req).Return(resp, nil)
	md := new(mockDirector)
	md.On("Direct", "get", "/", []byte(nil)).Return(req, nil)
	dc := &DirectedClient{mc, md}
	r, err := dc.IssueRequest("get", "/", nil)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, r.StatusCode)
}
