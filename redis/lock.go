package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func TryLock(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	ok, err := Client().SetNX(ctx, namespacedKey(key), value, ttl).Result()
	return ok, err
}

func Unlock(ctx context.Context, key string) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return Client().Del(ctx, namespacedKey(key)).Err()
}

var delIfEqual = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end`)

func SaferUnlock(ctx context.Context, key, value string) (bool, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	res, err := delIfEqual.Run(ctx, Client(), []string{namespacedKey(key)}, value).Result()
	if err != nil {
		return false, err
	}
	n, _ := res.(int64)
	return n == 1, nil
}
