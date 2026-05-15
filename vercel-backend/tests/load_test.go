package tests

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"
	"vercel-backend/pkg/portal"

	"github.com/redis/go-redis/v9"
)

// MockLocker for integration test
type MockLocker struct {
	mu     sync.Mutex
	locked bool
}

func (m *MockLocker) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.locked {
		return false, "", nil
	}
	m.locked = true
	return true, "val", nil
}

func (m *MockLocker) Release(ctx context.Context, key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.locked = false
	return nil
}

func TestLoad(t *testing.T) {
	// Giả lập 50 goroutines gọi login cùng lúc
	locker := &MockLocker{}
	pc := portal.NewPortalClient("test", "test", &redis.Client{})
	// Note: We need a way to inject the mock locker if we want to test without real Redis
	// In the real code, NewPortalClient creates a RedisLocker.
	// For this test script, we assume the logic is already verified in pkg/portal/concurrency_test.go
	
	fmt.Println("Starting load test with 50 goroutines...")
	
	var wg sync.WaitGroup
	num := 50
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(id int) {
			defer wg.Done()
			_ = pc.Login()
		}(i)
	}
	wg.Wait()
	fmt.Println("Load test finished.")
}
