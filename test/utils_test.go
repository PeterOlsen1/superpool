package superpool_test

import (
	"testing"

	sp "github.com/PeterOlsen1/superpool"
)

func setup(t *testing.T) *sp.Pool[int] {
	handler := func(i int) error {
		count += 1
		return nil
	}

	p, err := sp.NewPool(10, 2, handler)
	if err != nil {
		t.Fatal("pool creation failed")
	}
	return p
}

func setupCustom(t *testing.T, cap uint32, workers uint16, handler sp.EventHandler[int]) *sp.Pool[int] {
	p, err := sp.NewPool(cap, workers, handler)
	if err != nil {
		t.Fatal("pool creation failed")
	}
	return p
}

func teardown(p *sp.Pool[int]) {
	p.Shutdown()
}
