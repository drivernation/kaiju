package http

import "net/http"

// Client is an interface used for issuing HTTP requests. It allows for a modular http client making testing easier.
// Due to the way Go interfaces are resolved, A net/http.Client satisfies this interface. This interface should be
// used anywhere an http client is needed to facilitate flexibility and ease of testing.
type Client interface {
	// Execute req and return the response. An optional error is returned if execution of the request failed for whatever reason.
	// Note that a non 200 status code will NOT NECESSARILY return an error.
	Do(req *http.Request) (*http.Response, error)
}

// DirectedClient is a wrapper around a client and a director for convenience.
type DirectedClient struct {
	Client
	Director
}

// IssueRequest directs a request and executes it, returning the response object and an optional error if something went wrong.
func (c *DirectedClient) IssueRequest(method, resource string, body []byte) (*http.Response, error) {
	req, err := c.Direct(method, resource, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
