package impl

import (
	"math"
	"time"
	"yadroTest/internal/event_manager"
)

type Desks struct {
	Desks          []Desk
	DesksAvailable int
}

func NewDesks(num int) Desks {
	return Desks{
		Desks:          make([]Desk, num),
		DesksAvailable: num,
	}
}

func (d Desks) Occupy(deskNum int, occupationTime time.Time) error {
	if d.Desks[deskNum].Occupied {
		return event_manager.DeskIsOccupied
	}
	d.Desks[deskNum].Occupied = true
	d.Desks[deskNum].OccupiedAt = occupationTime
	d.DesksAvailable--
	return nil
}

func (d Desks) Free(deskNum int, currentTime time.Time) error {
	if !d.Desks[deskNum].Occupied {
		return event_manager.DeskAlreadyFree
	}
	d.Desks[deskNum].Occupied = false
	d.Desks[deskNum].Account(currentTime)
	d.DesksAvailable++
	return nil
}

func (d Desks) FindAvailable() (int, error) {
	if d.DesksAvailable == 0 {
		return -1, event_manager.DeskIsOccupied
	}
	for i := 0; i < len(d.Desks); i++ {
		if !d.Desks[i].Occupied {
			return i, nil
		}
	}
	return -1, event_manager.DeskIsOccupied
}

func (d Desks) AreAnyFree() bool {
	if d.DesksAvailable == 0 {
		return true
	}
	return false
}

type Desk struct {
	OccupiedAt time.Time
	Occupied   bool
	Revenue    int
	Price      int
}

func (d Desk) Account(currentTime time.Time) {
	timeDiff := currentTime.Sub(d.OccupiedAt).Hours()
	d.Revenue += d.Price * int(math.Ceil(timeDiff))
}
