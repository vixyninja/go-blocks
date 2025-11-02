package redis

import (
	"context"
	"fmt"
)

func NextID(ctx context.Context, key string) (int64, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return Client().Incr(ctx, namespacedKey(key)).Result()
}

func NextBatch(ctx context.Context, key string, n int64) (start, end int64, err error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	if n <= 0 {
		return 0, 0, fmt.Errorf("[pkg.redis] idalloc: n must be > 0")
	}
	end, err = Client().IncrBy(ctx, namespacedKey(key), n).Result()
	if err != nil {
		return 0, 0, err
	}
	start = end - n + 1
	return
}

func NextPrefixed(ctx context.Context, key, prefix string, pad int) (string, error) {
	id, err := NextID(ctx, key)
	if err != nil {
		return "", err
	}
	if pad < 1 {
		pad = 1
	}
	return fmt.Sprintf("%s%0*d", prefix, pad, id), nil
}
