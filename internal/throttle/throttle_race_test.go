package throttle_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

// TestReadyConcurrentSafe verifies that concurrent calls to Ready do not race.
func TestReadyConcurrentSafe(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			th.Ready()
		}()
	}
	wg.Wait()
}

// TestResetConcurrentSafe verifies that Reset and Ready may be called concurrently.
func TestResetConcurrentSafe(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			th.Ready()
		}()
		go func() {
			defer wg.Done()
			th.Reset()
		}()
	}
	wg.Wait()
}
