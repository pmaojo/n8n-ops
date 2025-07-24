package client

import "net/http"

// HTTPDoer represents the ability to execute an HTTP request.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPDoer wraps http.Client to satisfy HTTPDoer.
type DefaultHTTPDoer struct{ *http.Client }
