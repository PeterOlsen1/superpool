package superpool_test

import (
	"testing"
)

var count = 0

func TestPool(t *testing.T) {
	p := setup()

	p.Add(1)
	teardown(p)

	if count != 1 {
		t.Errorf("count (%d) does not match expected 1", count)
	}
}
