package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type httpClient struct {
	client *http.Client
	config *Config
}

func NewHTTPClient(config *Config) HTTPClient {
	if config == nil {
		config = DefaultConfig()
	}

	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &httpClient{
		client: client,
		config: config,
	}
}

func (c *httpClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	c.setHeaders(req, headers)
	return c.client.Do(req)
}

func (c *httpClient) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	c.setHeaders(req, headers)
	return c.client.Do(req)
}

func (c *httpClient) Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}

	c.setHeaders(req, headers)
	return c.client.Do(req)
}

func (c *httpClient) Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	c.setHeaders(req, headers)
	return c.client.Do(req)
}

func (c *httpClient) setHeaders(req *http.Request, customHeaders map[string]string) {
	for key, value := range c.config.DefaultHeaders {
		req.Header.Set(key, value)
	}

	if req.Header.Get("User-Agent") == "" && c.config.UserAgent != "" {
		req.Header.Set("User-Agent", c.config.UserAgent)
	}

	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}
}
