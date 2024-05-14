package src

import (
	"math"
	"time"
)

type DeskStorage struct {
	Desks          []Desk
	OccupiedDesks  map[string]int
	AvailableDesks Bitmask
	Price          int
}

func NewDesks(num, price int) *DeskStorage {
	return &DeskStorage{
		Desks:          make([]Desk, num),
		OccupiedDesks:  make(map[string]int),
		AvailableDesks: NewBitmask(num),
		Price:          price,
	}
}

func (d *DeskStorage) Occupy(clientName string, deskNum int, occupationTime time.Time) error {
	if d.AvailableDesks.Get(deskNum) {
		return DeskIsOccupied
	}

	if _, ok := d.OccupiedDesks[clientName]; ok {
		if err := d.Free(clientName, occupationTime); err != nil {
			return err
		}
	}
	d.OccupiedDesks[clientName] = deskNum
	d.Desks[deskNum].OccupiedAt = occupationTime
	d.AvailableDesks.Set(deskNum)
	return nil
}

func (d *DeskStorage) Free(name string, currentTime time.Time) error {
	deskNum, ok := d.OccupiedDesks[name]
	if !ok {
		return nil
	}

	delete(d.OccupiedDesks, name)
	d.AvailableDesks.Clear(deskNum)
	d.Desks[deskNum].Account(d.Price, currentTime)
	return nil
}

func (d *DeskStorage) FindAvailable() (int, bool) {
	return d.AvailableDesks.GetAvailable()
}

type Desk struct {
	OccupiedAt     time.Time
	OccupationTime time.Time // Should default to 00:00
	Revenue        int
}

func (d *Desk) Account(price int, currentTime time.Time) {
	timeDiff := currentTime.Sub(d.OccupiedAt)
	d.Revenue += price * int(math.Ceil(timeDiff.Hours()))
	d.OccupationTime = d.OccupationTime.Add(timeDiff)
}
