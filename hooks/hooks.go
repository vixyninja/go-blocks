package hooks

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vixyninja/go-blocks/logx"
)

type HookCallback func(ctx context.Context)

type Hook struct {
	timeout   time.Duration
	callbacks []HookCallback
	mu        sync.Mutex
}

func New(tineout time.Duration) *Hook {
	return &Hook{
		timeout:   tineout,
		callbacks: make([]HookCallback, 0),
	}
}

func (h *Hook) Add(cb HookCallback) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.callbacks = append(h.callbacks, cb)
}

func (h *Hook) Wait() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	h.mu.Lock()
	defer h.mu.Unlock()
	for _, cb := range h.callbacks {
		cb(ctx)
	}

	logx.NewStdLogger().Info(context.Background(), "[pkg.hooks.Wait] received signal %s, executing shutdown hook", sig.String())
}
