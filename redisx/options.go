package redis

import "time"

type Options struct {
	Addr              string
	Password          string
	DB                int
	UseTLS            bool
	PoolSize          int           // default 10*CPU
	MinIdleConns      int           // default 0
	DialTimeout       time.Duration // default 5s
	ReadTimeout       time.Duration // default 3s
	WriteTimeout      time.Duration // default 3s
	MaxRetries        int           // default 3
	MinRetryBackoff   time.Duration // default 8ms
	MaxRetryBackoff   time.Duration // default 512ms
	Namespace         string        // Namespacing key to avoid collisions with other projects
	DefaultTTL        time.Duration // Default TTL for cache (can be 0 = no expire)
	DefaultCmdTimeout time.Duration // Default context deadline for each command
}

func (o *Options) withDefaults() *Options {
	if o.DialTimeout == 0 {
		o.DialTimeout = 5 * time.Second
	}
	if o.ReadTimeout == 0 {
		o.ReadTimeout = 3 * time.Second
	}
	if o.WriteTimeout == 0 {
		o.WriteTimeout = 3 * time.Second
	}
	if o.MaxRetries == 0 {
		o.MaxRetries = 3
	}
	if o.MinRetryBackoff == 0 {
		o.MinRetryBackoff = 8 * time.Millisecond
	}
	if o.MaxRetryBackoff == 0 {
		o.MaxRetryBackoff = 512 * time.Millisecond
	}
	if o.DefaultCmdTimeout == 0 {
		o.DefaultCmdTimeout = 1 * time.Second
	}
	return o
}
