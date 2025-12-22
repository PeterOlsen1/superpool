package superpool_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestReturnPool(t *testing.T) {
	p := setupReturnPool(t)

	p.Add(1)
	teardownReturnPool(p)

	if count != 1 {
		t.Errorf("count (%d) does not match expected 1", count)
	}
}

func TestReturnPoolError(t *testing.T) {
	handler := func(i int) (int, error) {
		return 0, fmt.Errorf("testing")
	}
	p := setupCustomReturnPool(t, 1, 1, handler)

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
	teardownReturnPool(p)
	wg.Wait()

	if errorCount != 2 {
		t.Errorf("the error channel was not signaled twice")
	}
}

func TestReturnPoolReturn(t *testing.T) {
	p := setupReturnPool(t)

	retVal := 0
	for c := range p.Add(1) {
		retVal = c
	}

	if retVal != 1 {
		t.Errorf("Pool did not return 1")
	}

	// calls the shutdown method
	teardownReturnPool(p)
}
