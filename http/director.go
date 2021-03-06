package http

import (
	"bytes"
	"io"
	"net/http"
	"path"
)

type Scheme string

const (
	SchemeHttp  = Scheme("http://")
	SchemeHttps = Scheme("https://")
)

type Director interface {
	Direct(method, resource string, body []byte) (*http.Request, error)
}

type roundRobin struct {
	hosts  []string
	scheme Scheme
}

func (r *roundRobin) Direct(method, resource string, body []byte) (*http.Request, error) {
	var next string
	next, r.hosts = r.hosts[0], r.hosts[1:]
	r.hosts = append(r.hosts, next)
	url := path.Join(string(r.scheme), next, resource)
	var br io.Reader = nil
	if body != nil {
		br = bytes.NewBuffer(body)
	}
	return http.NewRequest(method, url, br)
}

func RoundRobin(hosts []string, scheme Scheme) Director {
	return &roundRobin{
		hosts:  hosts,
		scheme: scheme,
	}
}
