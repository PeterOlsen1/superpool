package benchmarks

import (
	"net/http"
	"testing"

	sp "github.com/PeterOlsen1/superpool"
)

var taskCount int = 100
var server *http.Server

func startLocal() {
	server = &http.Server{Addr: "127.0.0.1:8000"}
	go func() {
		server.ListenAndServe()
	}()
}

func stopLocal() {
	if server != nil {
		server.Shutdown(nil)
	}
}

func dummyTask(i int) error {
	http.Get("127.0.0.1:8000")
	return nil
}

// BenchmarkPool-16           18999             60530 ns/op           20815 B/op        300 allocs/op
func BenchmarkPool(b *testing.B) {
	startLocal()
	defer stopLocal()

	numWorkers := uint16(50)
	poolSize := uint32(1000)
	pool, _ := sp.NewPool(poolSize, numWorkers, dummyTask)
	defer pool.Shutdown()

	for b.Loop() {
		for range taskCount {
			pool.Add(1)
		}
		pool.Wait()
	}
}

// BenchmarkGoroutines-16    	   19046	     61192 ns/op	   22531 B/op	     401 allocs/op
func BenchmarkGoroutines(b *testing.B) {
	startLocal()
	defer stopLocal()

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

// BenchmarkSequential-16    	   77817	     14602 ns/op	   20801 B/op	     300 allocs/op
func BenchmarkSequential(b *testing.B) {
	startLocal()
	defer stopLocal()

	for b.Loop() {
		for range taskCount {
			dummyTask(1)
		}
	}
}
