package superpool_test

import (
	sp "github.com/PeterOlsen1/superpool"
)

func setup() *sp.Pool[int] {
	handler := func(i int) error {
		count += 1
		return nil
	}

	p := sp.NewPool(10, 2, handler)
	return p
}

func teardown(p *sp.Pool[int]) {
	p.Close()
}
