package src

import "math"

type Bitmask struct {
	bitmask []uint64
	size    int
}

func NewBitmask(size int) Bitmask {
	return Bitmask{
		bitmask: make([]uint64, (size+63)/64),
		size:    size,
	}
}

func (b Bitmask) Get(id int) bool {
	return (b.bitmask[id/64] & (1 << (id % 64))) > 0
}

func (b Bitmask) Set(id int) {
	b.bitmask[id/64] |= 1 << (id % 64)
}

func (b Bitmask) Clear(id int) {
	b.bitmask[id/64] &= ^(1 << (id % 64))
}

func (b Bitmask) GetAvailable() (int, bool) {
	for i, val := range b.bitmask {
		if val != math.MaxUint64 {
			upperBound := 64
			if i == len(b.bitmask)-1 {
				upperBound = b.size % 64
			}
			for j := 0; j < upperBound; j++ {
				if val&(1<<(j)) == 0 {
					return i*64 + j, true
				}
			}
		}
	}
	return 0, false
}
