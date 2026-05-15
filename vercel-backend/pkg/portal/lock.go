package portal

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Locker defines the interface for distributed locking
type Locker interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, string, error)
	Release(ctx context.Context, key string, value string) error
}

// RedisLocker implements Locker using Redis
type RedisLocker struct {
	client *redis.Client
}

// NewRedisLocker creates a new RedisLocker
func NewRedisLocker(client *redis.Client) *RedisLocker {
	return &RedisLocker{client: client}
}

// Acquire tries to acquire a lock with a unique value and TTL
func (l *RedisLocker) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, string, error) {
	if l.client == nil {
		return true, "", nil // No redis, no lock
	}
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	ok, err := l.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, "", err
	}
	return ok, value, nil
}

var releaseScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`)

// Release safely releases the lock using a Lua script
func (l *RedisLocker) Release(ctx context.Context, key string, value string) error {
	if l.client == nil {
		return nil
	}
	return releaseScript.Run(ctx, l.client, []string{key}, value).Err()
}
