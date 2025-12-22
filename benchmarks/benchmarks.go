package benchmarks

import (
	"testing"
	"time"

	sp "github.com/PeterOlsen1/superpool"
)

var taskCount int = 10000

func dummyTask(i int) error {
	// Simulate work (e.g., 1ms sleep)
	time.Sleep(1 * time.Millisecond)
	return nil
}

func BenchmarkSuperpool(b *testing.B) {
	poolSize := uint16(50)

	for n := 0; n < b.N; n++ {
		pool, _ := sp.NewPool(1000, poolSize, dummyTask)
		for i := 0; i < taskCount; i++ {
			pool.Add(1)
		}
		pool.Wait()
	}
}

func BenchmarkGoroutines(b *testing.B) {
	for n := 0; n < b.N; n++ {
		done := make(chan struct{}, taskCount)
		for i := 0; i < taskCount; i++ {
			go func() {
				dummyTask(1)
				done <- struct{}{}
			}()
		}
		for i := 0; i < taskCount; i++ {
			<-done
		}
	}
}

func BenchmarkSequential(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < taskCount; i++ {
			dummyTask(1)
		}
	}
}
