package http

import (
	"context"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
	Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
}

type Config struct {
	Timeout        time.Duration
	UserAgent      string
	DefaultHeaders map[string]string
}

func DefaultConfig() *Config {
	return &Config{
		Timeout:        30 * time.Second,
		UserAgent:      "Gaia-Client/1.0",
		DefaultHeaders: make(map[string]string),
	}
}
