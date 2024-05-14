package src

import "math"

type Bitmask struct {
	bitmask []uint64
	size    int
}

func NewBitmask(size int) Bitmask {
	return Bitmask{
		bitmask: make([]uint64, (size+7)/8),
		size:    size,
	}
}

func (b Bitmask) Get(id int) bool {
	return (b.bitmask[id/8] & (1 << (id % 8))) > 0
}

func (b Bitmask) Set(id int) {
	b.bitmask[id/8] |= 1 << (id % 8)
}

func (b Bitmask) Clear(id int) {
	b.bitmask[id/8] &= ^(1 << (id % 8))
}

func (b Bitmask) GetAvailable() (int, bool) {
	for i, val := range b.bitmask {
		if val != math.MaxUint64 {
			upperBound := 8
			if i == len(b.bitmask)-1 {
				upperBound = b.size % 8
			}
			for j := 0; j < upperBound; j++ {
				if val&(1<<(j)) == 0 {
					return i*8 + j, true
				}
			}
		}
	}
	return 0, false
}
