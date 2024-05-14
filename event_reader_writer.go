package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"yadroTest/src"
)

type Config struct {
	OpeningTime time.Time
	ClosingTime time.Time
	DesksNumber int
	Price       int
}

type EventReaderWriter struct {
	EventManager        src.EventManager
	OpeningTime         time.Time
	ClosingTime         time.Time
	MostRecentEventTime time.Time
	DesksNumber         int
	HourlyPrice         int
}

func ReadInitData(r io.Reader) (*Config, error) {
	var numDesks int
	_, err := fmt.Fscanln(r, &numDesks)
	if err != nil {
		return nil, err
	}

	var openingTimeStr, closingTimeStr string
	_, err = fmt.Fscanln(r, &openingTimeStr, &closingTimeStr)
	if err != nil {
		return nil, err
	}
	openingTime, err := time.Parse(HHMM24H, openingTimeStr)
	if err != nil {
		return nil, err
	}
	closingTime, err := time.Parse(HHMM24H, closingTimeStr)

	var price int
	_, err = fmt.Fscanln(r, &price)
	if err != nil {
		return nil, err
	}

	return &Config{
		OpeningTime: openingTime,
		ClosingTime: closingTime,
		DesksNumber: numDesks,
		Price:       price,
	}, nil
}

func (e *EventReaderWriter) ReadEvent(eventStr string) (string, error) {
	eventInfo := strings.Fields(eventStr)
	eventTime, err := time.Parse(HHMM24H, eventInfo[0])
	if err != nil {
		return "", err
	}
	if eventTime.Before(e.OpeningTime) || eventTime.After(e.ClosingTime) {
		sideEffect := fmt.Sprintf(" %d %s\n", 13, src.NotOpenYet)
		return eventTime.Format(HHMM24H) + sideEffect, nil
	}
	if eventTime.Before(e.MostRecentEventTime) {
		return "", src.EventFormatError
	}
	eventID, err := strconv.Atoi(eventInfo[1])
	if err != nil {
		return "", src.EventFormatError
	}
	clientName := eventInfo[2]

	var deskNum int = -1
	if len(eventInfo) == 4 {
		deskNum, err = strconv.Atoi(eventInfo[3])
		if err != nil {
			return "", err
		}
		if deskNum > e.DesksNumber {
			return "", src.EventFormatError
		}
	}

	var sideEffect string
	switch eventID {
	case 1:
		sideEffect = e.ClientArrived(clientName)
	case 2:
		sideEffect = e.ClientSatAtTheDesk(deskNum-1, clientName, eventTime)
	case 3:
		sideEffect = e.ClientAwaits(clientName, eventTime)
	case 4:
		sideEffect = e.ClientLeaves(clientName, eventTime)
	default:
		return "", src.EventFormatError
	}
	if len(sideEffect) == 0 {
		return "", err
	}
	return eventTime.Format(HHMM24H) + " " + sideEffect, err
}

func Handle(r *bufio.Reader, w io.Writer) error {
	config, err := ReadInitData(r)
	if err != nil {
		return err
	}
	eventManager := src.EventManager{
		ClientPool:  src.NewClientPool(),
		DeskStorage: src.NewDesks(config.DesksNumber, config.Price),
		ClientQueue: src.NewClientListQueue(),
	}
	eventReaderWriter := &EventReaderWriter{
		EventManager:        eventManager,
		OpeningTime:         config.OpeningTime,
		ClosingTime:         config.ClosingTime,
		MostRecentEventTime: config.OpeningTime,
		DesksNumber:         config.DesksNumber,
		HourlyPrice:         config.Price,
	}

	// Write opening time
	fmt.Fprint(w, config.OpeningTime.Format(HHMM24H)+"\n")

	// Handle all events
	for eventStr, err := r.ReadString('\n'); err == nil; eventStr, err = r.ReadString('\n') {
		sideEffectStr, err := eventReaderWriter.ReadEvent(eventStr)
		fmt.Fprint(w, eventStr)
		if errors.Is(err, src.EventFormatError) {
			break
		}
		if len(sideEffectStr) > 0 {
			fmt.Fprint(w, sideEffectStr)
		}
	}
	// Kick out all customers
	for name, _ := range eventReaderWriter.EventManager.ClientPool.Pool {
		eventReaderWriter.EventManager.DeskStorage.Free(name, config.ClosingTime)
		kickOutEvent := fmt.Sprintf("%s %d %s\n", config.ClosingTime.Format(HHMM24H), 11, name)
		fmt.Fprint(w, kickOutEvent)
	}
	// Write closing time
	fmt.Fprint(w, config.ClosingTime.Format(HHMM24H)+"\n")

	// Write all desks' statistics
	for i, desk := range eventReaderWriter.EventManager.DeskStorage.Desks {
		deskInfo := fmt.Sprintf("%d %d %s\n", i+1, desk.Revenue, desk.OccupationTime.Format(HHMM24H))
		fmt.Fprint(w, deskInfo)
	}
	return nil
}

func (e *EventReaderWriter) ClientArrived(name string) string {
	if err := e.EventManager.ClientArrived(name); err != nil {
		return fmt.Sprintf("%d %v\n", 13, err)
	}
	return ""
}

func (e *EventReaderWriter) ClientSatAtTheDesk(deskNum int, name string, currentTime time.Time) string {
	if err := e.EventManager.ClientSatAtTheDesk(deskNum, name, currentTime); err != nil {
		return fmt.Sprintf("%d %v\n", 13, err)
	}
	return ""
}

func (e *EventReaderWriter) ClientAwaits(name string, currentTime time.Time) string {
	if err := e.EventManager.ClientAwaits(name); errors.Is(err, src.QueueIsFull) {
		return fmt.Sprintf("%d %s\n", 11, name)
	} else if errors.Is(err, src.ICanWaitNoLonger) {
		return fmt.Sprintf("%d %s\n", 13, err)
	}
	return ""
}

func (e *EventReaderWriter) ClientLeaves(name string, currentTime time.Time) string {
	if err := e.EventManager.ClientLeaves(name, currentTime); err != nil {
		return fmt.Sprintf("%d %v\n", 13, err)
	}
	deskNum, ok := e.EventManager.DeskStorage.FindAvailable()
	if !ok {
		return ""
	}
	awaitingClient, ok := e.EventManager.ClientQueue.Dequeue()
	if !ok {
		return ""
	}
	if err := e.EventManager.ClientSatAtTheDesk(deskNum, awaitingClient, currentTime); err != nil {
		return fmt.Sprintf("%d %v\n", 13, err)
	}
	return fmt.Sprintf("%d %s %d\n", 12, awaitingClient, deskNum+1)
}

const HHMM24H = "15:04"
