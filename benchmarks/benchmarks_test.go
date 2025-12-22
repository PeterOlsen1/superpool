package benchmarks

import (
	"testing"
	"time"

	sp "github.com/PeterOlsen1/superpool"
)

var taskCount int = 10000

func dummyTask(i int) error {
	time.Sleep(1 * time.Millisecond)
	return nil
}

// BenchmarkPool-16    	      10	 111145265 ns/op	   93129 B/op	     544 allocs/op
func BenchmarkPool(b *testing.B) {
	poolSize := uint16(50)

	for b.Loop() {
		pool, _ := sp.NewPool(1000, poolSize, dummyTask)
		for range taskCount {
			pool.Add(1)
		}
		pool.Shutdown()
	}
}

// BenchmarkGoroutines-16    	     196	   6082681 ns/op	 1143147 B/op	   20061 allocs/op
func BenchmarkGoroutines(b *testing.B) {
	for b.Loop() {
		done := make(chan struct{}, taskCount)
		for range taskCount {
			go func() {
				dummyTask(1)
				done <- struct{}{}
			}()
		}
		for range taskCount {
			<-done
		}
	}
}

// BenchmarkSequential-16    	       1	10705755996 ns/op	     176 B/op	       5 allocs/op
func BenchmarkSequential(b *testing.B) {
	for b.Loop() {
		for range taskCount {
			dummyTask(1)
		}
	}
}
