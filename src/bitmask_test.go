package src

import "testing"

func TestBitmask(t *testing.T) {
	bitmask := NewBitmask(10)
	bitmask.Set(5)
	bitmask.Set(9)

	for i := 0; i < 10; i++ {
		if bitmask.Get(i) != (i == 5 || i == 9) {
			t.Fatal("expected bitmask to not be set")
		}
	}
	bitmask.Clear(9)
	for i := 0; i < 10; i++ {
		if bitmask.Get(i) != (i == 5) {
			t.Fatal("expected bitmask to not be set")
		}
	}

	for val, ok := bitmask.GetAvailable(); ok; val, ok = bitmask.GetAvailable() {
		bitmask.Set(val)
	}
	for i := 0; i < 10; i++ {
		if !bitmask.Get(i) {
			t.Fatal("expected bitmask to be set")
		}
	}
}
