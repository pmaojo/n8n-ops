package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	headerAPIKey      = "X-N8N-API-KEY"
	headerAccept      = "Accept"
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"

	defaultTimeout   = 30 * time.Second
	handshakeTimeout = 10 * time.Second
	headerTimeout    = 15 * time.Second
	idleTimeout      = 90 * time.Second
)

// StandardHTTPTransport implements HTTPTransport using an HTTPDoer
type StandardHTTPTransport struct {
	client  HTTPDoer
	apiKey  string
	baseURL string
}

// NewStandardHTTPTransport creates a new HTTP transport with optimal configuration
func NewStandardHTTPTransport(baseURL, apiKey string, client HTTPDoer) *StandardHTTPTransport {
	if client == nil {
		client = &DefaultHTTPDoer{defaultHTTPClient()}
	}

	return &StandardHTTPTransport{
		client:  client,
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// defaultHTTPClient creates an optimally configured HTTP client
func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			TLSHandshakeTimeout:   handshakeTimeout,
			ResponseHeaderTimeout: headerTimeout,
			IdleConnTimeout:       idleTimeout,
			MaxIdleConnsPerHost:   10,
			DisableKeepAlives:     false,
		},
	}
}

// Do executes the HTTP request with proper error handling
func (t *StandardHTTPTransport) Do(req *Request) (*Response, error) {
	return t.doRequest(context.Background(), req)
}

func (t *StandardHTTPTransport) doRequest(ctx context.Context, req *Request) (*Response, error) {
	var body io.Reader
	if len(req.Body) > 0 {
		body = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, t.baseURL+req.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	t.addDefaultHeaders(httpReq)

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	httpResp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	headers := make(map[string]string)
	for key, values := range httpResp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    headers,
		Body:       respBody,
	}, nil
}

// addDefaultHeaders adds standard headers to every request
func (t *StandardHTTPTransport) addDefaultHeaders(req *http.Request) {
	req.Header.Set(headerAPIKey, t.apiKey)
	req.Header.Set(headerAccept, contentTypeJSON)

	// Set Content-Type for requests with body
	if req.Body != nil {
		req.Header.Set(headerContentType, contentTypeJSON)
	}
}

// ContextualTransport wraps StandardHTTPTransport with context support
type ContextualTransport struct {
	transport *StandardHTTPTransport
}

// NewContextualTransport creates a transport that supports context cancellation
func NewContextualTransport(baseURL, apiKey string, client HTTPDoer) *ContextualTransport {
	return &ContextualTransport{
		transport: NewStandardHTTPTransport(baseURL, apiKey, client),
	}
}

// DoWithContext executes request with context support
func (t *ContextualTransport) DoWithContext(ctx context.Context, req *Request) (*Response, error) {
	return t.transport.doRequest(ctx, req)
}
