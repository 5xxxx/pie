package pie

import (
        "context"
        "sync"
        "sync/atomic"
        "testing"
        "time"
)

func TestSingleFlightDo(t *testing.T) {
        sf := NewSingleFlight()
        var counter int32
        var wg sync.WaitGroup
        wg.Add(10)

        for i := 0; i < 10; i++ {
                go func() {
                        defer wg.Done()
                        val, err := sf.Do("key", func() ([]byte, error) {
                                time.Sleep(5 * time.Millisecond)
                                atomic.AddInt32(&counter, 1)
                                return []byte("value"), nil
                        })
                        if err != nil {
                                t.Errorf("Do returned error: %v", err)
                        }
                        if string(val) != "value" {
                                t.Errorf("unexpected value: %s", val)
                        }
                }()
        }
        wg.Wait()

        if counter != 1 {
                t.Fatalf("function should execute once, executed %d times", counter)
        }
}

func TestCacheWithSingleFlightDelegation(t *testing.T) {
        mock := newMockCache()
        cache := NewCacheWithSingleFlight(mock)
        ctx := context.Background()

        if err := cache.Set(ctx, "key", []byte("value"), time.Second); err != nil {
                t.Fatalf("Set error: %v", err)
        }

        var wg sync.WaitGroup
        wg.Add(2)
        go func() {
                defer wg.Done()
                if _, err := cache.Get(ctx, "key"); err != nil {
                        t.Errorf("Get error: %v", err)
                }
        }()
        go func() {
                defer wg.Done()
                if _, err := cache.Get(ctx, "key"); err != nil {
                        t.Errorf("Get error: %v", err)
                }
        }()
        wg.Wait()

        if err := cache.Delete(ctx, "key"); err != nil {
                t.Fatalf("Delete error: %v", err)
        }
        if err := cache.DeleteByPattern(ctx, "prefix"); err != nil {
                t.Fatalf("DeleteByPattern error: %v", err)
        }
        if err := cache.DeleteByTags(ctx, "a"); err != nil {
                t.Fatalf("DeleteByTags error: %v", err)
        }
        if exists, err := cache.Exists(ctx, "key"); err != nil || exists {
                t.Fatalf("Exists returned (%v, %v)", err, exists)
        }
        if err := cache.Clear(ctx); err != nil {
                t.Fatalf("Clear error: %v", err)
        }
        if cache.Stats() == nil {
                t.Fatalf("Stats returned nil")
        }
}

