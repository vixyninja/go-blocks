package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
	opts   *Options
)

func Init(o Options) error {
	var initErr error
	once.Do(func() {
		opts = o.withDefaults()

		ro := &redis.Options{
			Addr:            opts.Addr,
			Password:        opts.Password,
			DB:              opts.DB,
			PoolSize:        opts.PoolSize,
			MinIdleConns:    opts.MinIdleConns,
			DialTimeout:     opts.DialTimeout,
			ReadTimeout:     opts.ReadTimeout,
			WriteTimeout:    opts.WriteTimeout,
			MaxRetries:      opts.MaxRetries,
			MinRetryBackoff: opts.MinRetryBackoff,
			MaxRetryBackoff: opts.MaxRetryBackoff,
		}
		if opts.UseTLS {
			ro.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		client = redis.NewClient(ro)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("[pkg.redis] ping failed: %w", err)
			return
		}
	})
	return initErr
}

func Client() *redis.Client {
	if client == nil {
		panic("[pkg.redis] not initialized, call redis.Init first")
	}
	return client
}

func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

func namespacedKey(key string) string {
	if opts == nil || opts.Namespace == "" {
		return key
	}
	return opts.Namespace + ":" + key
}

func withTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	d := opts.DefaultCmdTimeout
	if d <= 0 {
		return parent, func() {}
	}
	return context.WithTimeout(parent, d)
}
