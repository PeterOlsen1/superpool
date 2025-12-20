package superpool_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var count = 0

func TestAdd(t *testing.T) {
	p := setup(t)

	p.Add(1)
	teardown(p)

	if count != 1 {
		t.Errorf("count (%d) does not match expected 1", count)
	}
}

func TestError(t *testing.T) {
	handler := func(i int) error {
		return fmt.Errorf("testing")
	}
	p := setupCustom(t, 1, 1, handler)

	var wg sync.WaitGroup
	wg.Add(1)

	p.Add(1)
	errorCount := 0
	go func() {
		defer wg.Done()
		for range p.Errors() {
			errorCount += 1
		}
	}()
	p.Add(1)

	time.Sleep(1 * time.Millisecond)
	// calls the shutdown method
	teardown(p)
	wg.Wait()

	if errorCount != 2 {
		t.Errorf("the error channel was not signaled twice")
	}
}

func TestUpdateEventHandler(t *testing.T) {
	p := setup(t)

	p.Add(1)
	time.Sleep(1 * time.Millisecond)
	if count != 1 {
		t.Errorf("count (%d) does not match expected 1", count)
	}

	p.UpdateEventHandler(func(i int) error {
		count += 2
		return nil
	})

	p.Add(1)
	time.Sleep(1 * time.Millisecond)
	if count != 3 {
		t.Errorf("count (%d) does not match expected 3", count)
	}
}
