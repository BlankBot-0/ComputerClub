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
		return "", src.NotOpenYet
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

	var sideEffect *SideEffect
	switch eventID {
	case 1:
		err = e.EventManager.ClientArrived(clientName)
	case 2:
		err = e.EventManager.ClientSatAtTheDesk(deskNum, clientName, eventTime)
	case 3:
		sideEffect = e.ClientAwaits(clientName, eventTime)
	case 4:
		sideEffect, err = e.ClientLeaves(clientName, eventTime)
	default:
		return "", src.EventFormatError
	}
	if sideEffect == nil {
		return "", err
	}
	return eventTime.String() + sideEffect.String(), err
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

	for eventStr, err := r.ReadString('\n'); err == nil; eventStr, err = r.ReadString('\n') {
		sideEffectStr, err := eventReaderWriter.ReadEvent(eventStr)
		fmt.Fprintln(w, eventStr)
		if errors.Is(err, src.EventFormatError) {
			break
		}
		if err != nil {
			fmt.Fprintln(w, SideEffect{
				SideEffectID: 13,
				ClientName:   "",
				DeskNumber:   0,
			})
		}
		if len(sideEffectStr) > 0 {
			fmt.Fprintln(w, sideEffectStr)
		}
	}
	return nil
}

func (e *EventReaderWriter) ClientAwaits(name string, currentTime time.Time) *SideEffect {
	if err := e.EventManager.ClientAwaits(name); errors.Is(err, src.QueueIsFull) {
		return &SideEffect{
			FormattedTime: currentTime.Format(HHMM24H),
			SideEffectID:  11,
			ClientName:    name,
			DeskNumber:    0,
		}
	}
	return nil
}

func (e *EventReaderWriter) ClientLeaves(name string, currentTime time.Time) (*SideEffect, error) {
	if err := e.EventManager.ClientLeaves(name, currentTime); err != nil {
		return nil, err
	}
	deskNum, ok := e.EventManager.DeskStorage.FindAvailable()
	if !ok {
		return nil, nil
	}
	awaitingClient, ok := e.EventManager.ClientQueue.Dequeue()
	if !ok {
		return nil, nil
	}
	if err := e.EventManager.ClientSatAtTheDesk(deskNum, awaitingClient, currentTime); err != nil {
		return nil, err
	}
	return &SideEffect{
		FormattedTime: currentTime.Format(HHMM24H),
		SideEffectID:  12,
		ClientName:    awaitingClient,
		DeskNumber:    deskNum,
	}, nil
}

type SideEffect struct {
	FormattedTime string
	SideEffectID  int
	ClientName    string
	DeskNumber    int
}

func (s SideEffect) String() string {
	return fmt.Sprintf("%s %d %s %d\n", s.FormattedTime, s.SideEffectID, s.ClientName, s.DeskNumber)
}

const HHMM24H = "15:04"
