package redis

import (
	"context"
	"encoding/json"
	"time"
)

func SetString(ctx context.Context, key, val string, ttl time.Duration) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return Client().Set(ctx, namespacedKey(key), val, ttl).Err()
}

func GetString(ctx context.Context, key string) (string, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return Client().Get(ctx, namespacedKey(key)).Result()
}

func Del(ctx context.Context, keys ...string) (int64, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	nsKeys := make([]string, len(keys))
	for i, k := range keys {
		nsKeys[i] = namespacedKey(k)
	}
	return Client().Del(ctx, nsKeys...).Result()
}

func SetJSON(ctx context.Context, key string, v any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if ttl == 0 && opts != nil && opts.DefaultTTL > 0 {
		ttl = opts.DefaultTTL
	}
	return SetString(ctx, key, string(b), ttl)
}

func GetJSON(ctx context.Context, key string, out any) (bool, error) {
	s, err := GetString(ctx, key)
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	}
	return true, json.Unmarshal([]byte(s), out)
}

func Incr(ctx context.Context, key string) (int64, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return Client().Incr(ctx, namespacedKey(key)).Result()
}

func Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return Client().Expire(ctx, namespacedKey(key), ttl).Result()
}

func HSet(ctx context.Context, key string, values map[string]any) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return Client().HSet(ctx, namespacedKey(key), values).Err()
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return Client().HGetAll(ctx, namespacedKey(key)).Result()
}
