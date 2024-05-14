package src

import "testing"

func TestBitmask(t *testing.T) {
	bitmask := NewBitmask(70)
	bitmask.Set(5)
	bitmask.Set(9)
	bitmask.Set(65)
	bitmask.Set(62)

	for i := 0; i < 70; i++ {
		if bitmask.Get(i) != (i == 5 || i == 9 || i == 65 || i == 62) {
			t.Fatal("expected bitmask to not be set")
		}
	}
	bitmask.Clear(9)
	bitmask.Clear(65)
	bitmask.Clear(62)
	for i := 0; i < 70; i++ {
		if bitmask.Get(i) != (i == 5) {
			t.Fatal("expected bitmask to not be set")
		}
	}

	for val, ok := bitmask.GetAvailable(); ok; val, ok = bitmask.GetAvailable() {
		bitmask.Set(val)
	}
	for i := 0; i < 70; i++ {
		if !bitmask.Get(i) {
			t.Fatal("expected bitmask to be set")
		}
	}
}
