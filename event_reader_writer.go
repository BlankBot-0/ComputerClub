package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
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

func ReadInitData(r *bufio.Reader, w io.Writer) (*Config, error) {
	var numDesks int
	numDesksStr, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	numDesks, err = strconv.Atoi(numDesksStr[:len(numDesksStr)-1])
	if err != nil {
		_, err := fmt.Fprint(w, numDesksStr)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	var openingTimeStr, closingTimeStr string
	workTimeStr, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	workTimeStrFields := strings.Fields(workTimeStr)
	openingTimeStr = workTimeStrFields[0]
	closingTimeStr = workTimeStrFields[1]

	isCorrectTimeFormat := regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`).MatchString
	if !isCorrectTimeFormat(openingTimeStr) || !isCorrectTimeFormat(closingTimeStr) {
		_, err := fmt.Fprint(w, workTimeStr)
		if err != nil {
			return nil, err
		}
		return nil, src.EventFormatError
	}
	openingTime, err := time.Parse(HHMM24H, openingTimeStr)
	if err != nil {
		return nil, err
	}
	closingTime, err := time.Parse(HHMM24H, closingTimeStr)
	if err != nil {
		return nil, err
	}

	var price int
	priceStr, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	price, err = strconv.Atoi(priceStr[:len(priceStr)-1])
	if err != nil {
		_, err := fmt.Fprint(w, priceStr)
		if err != nil {
			return nil, err
		}
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
	if len(eventInfo) < 3 || len(eventInfo) > 4 {
		return "", src.EventFormatError
	}
	if !regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`).MatchString(eventInfo[0]) {
		return "", src.EventFormatError
	}

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
	e.MostRecentEventTime = eventTime
	eventID, err := strconv.Atoi(eventInfo[1])
	if err != nil {
		return "", src.EventFormatError
	}
	clientName := eventInfo[2]
	if !regexp.MustCompile(`^[A-Za-z0-9\-_]+$`).MatchString(clientName) {
		return "", src.EventFormatError
	}

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

func Handle(r *bufio.Reader, wSource io.Writer) error {
	w := bufio.NewWriter(wSource)
	config, err := ReadInitData(r, w)
	if err != nil {
		err := w.Flush()
		if err != nil {
			return err
		}
		return nil
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
	_, err = fmt.Fprint(w, config.OpeningTime.Format(HHMM24H)+"\n")
	if err != nil {
		return err
	}

	// Handle all events
	for eventStr, err := r.ReadString('\n'); err == nil; eventStr, err = r.ReadString('\n') {
		sideEffectStr, eventErr := eventReaderWriter.ReadEvent(eventStr)
		_, err = fmt.Fprint(w, eventStr)
		if err != nil {
			return err
		}
		if errors.Is(eventErr, src.EventFormatError) {
			w = bufio.NewWriter(wSource)
			_, err := fmt.Fprint(w, eventStr)
			if err != nil {
				return err
			}
			err = w.Flush()
			if err != nil {
				return err
			}
			return nil
		}
		if len(sideEffectStr) > 0 {
			_, err := fmt.Fprint(w, sideEffectStr)
			if err != nil {
				return err
			}
		}
	}
	// Kick out all customers
	sortedNames := make([]string, len(eventReaderWriter.EventManager.ClientPool.Pool))
	i := 0
	for name := range eventReaderWriter.EventManager.ClientPool.Pool {
		sortedNames[i] = name
		i++
	}
	sort.Strings(sortedNames)
	for _, name := range sortedNames {
		err := eventReaderWriter.EventManager.DeskStorage.Free(name, config.ClosingTime)
		if err != nil {
			return err
		}
		kickOutEvent := fmt.Sprintf("%s %d %s\n", config.ClosingTime.Format(HHMM24H), 11, name)
		_, err = fmt.Fprint(w, kickOutEvent)
		if err != nil {
			return err
		}
	}

	// Write closing time
	_, err = fmt.Fprint(w, config.ClosingTime.Format(HHMM24H)+"\n")
	if err != nil {
		return err
	}

	// Write all desks' statistics
	for i, desk := range eventReaderWriter.EventManager.DeskStorage.Desks {
		deskInfo := fmt.Sprintf("%d %d %s\n", i+1, desk.Revenue, desk.OccupationTime.Format(HHMM24H))
		_, err := fmt.Fprint(w, deskInfo)
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	if err != nil {
		return err
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
		err := e.EventManager.ClientLeaves(name, currentTime)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("%d %s\n", 11, name)
	} else if err != nil {
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
