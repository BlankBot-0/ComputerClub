package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"yadroTest/internal/event_manager"
)

type (
	EventManager interface {
		ClientArrived(name event_manager.ClientName) error
		ClientSatAtTheDesk(deskNum int, name event_manager.ClientName) error
		ClientAwaits(name event_manager.ClientName) error
		ClientLeaves(name event_manager.ClientName, currentTime time.Time) error
	}
)

type EventReaderWriter struct {
	EventProcessor      EventManager
	OpeningTime         time.Time
	ClosingTime         time.Time
	MostRecentEventTime time.Time
	DesksNumber         int
	HourlyPrice         int
}

func (e *EventReaderWriter) ReadInitData(r io.Reader) error {
	var numDesks int
	_, err := fmt.Fscanln(r, &numDesks)
	if err != nil {
		return err
	}

	var openingTime, cLosingTime time.Time
	_, err = fmt.Fscanln(r, &openingTime, &cLosingTime)
	if err != nil {
		return err
	}

	var hourlyPrice int
	_, err = fmt.Fscanln(r, &hourlyPrice)
	if err != nil {
		return err
	}

	e.DesksNumber = numDesks
	e.OpeningTime = openingTime
	e.ClosingTime = cLosingTime
	e.HourlyPrice = hourlyPrice
	return nil
}

func (e *EventReaderWriter) ReadEvent(w io.Reader) error {
	var eventString string
	_, err := fmt.Fscanln(w, &eventString)
	if err != nil {
		return err
	}

	eventInfo := strings.Fields(eventString)
	eventTime, err := time.Parse(HHMM24H, eventInfo[0])
	if err != nil {
		return err
	}
	if eventTime.Before(e.OpeningTime) || eventTime.After(e.ClosingTime) {
		return event_manager.NotOpenYet
	}
	if eventTime.Before(e.MostRecentEventTime) {
		return event_manager.EventFormatError
	}
	eventID, err := strconv.Atoi(eventInfo[1])
	if err != nil {
		return event_manager.EventFormatError
	}
	//clientName := eventInfo[2]

	if len(eventInfo) == 4 {
		tableNumber, err := strconv.Atoi(eventInfo[3])
		if err != nil {
			return err
		}
		if tableNumber > e.DesksNumber {
			return event_manager.EventFormatError
		}
	}

	switch eventID {
	case 1:

	}
	return nil
}

func (e *EventReaderWriter) ClientArrived(name event_manager.ClientName) error {
	return e.EventProcessor.ClientArrived(name)
}

func (e *EventReaderWriter) ClientSatAtTheDesk(deskNum int, name event_manager.ClientName) error {
	return e.EventProcessor.ClientSatAtTheDesk(deskNum, name)
}

func (e *EventReaderWriter) ClientAwaits(name event_manager.ClientName) error {
	return e.EventProcessor.ClientAwaits(name)
}

func (e *EventReaderWriter) ClientLeaves(name event_manager.ClientName, currentTime time.Time) error {
	return e.EventProcessor.ClientLeaves(name, currentTime)
}

const HHMM24H = "15:04"
