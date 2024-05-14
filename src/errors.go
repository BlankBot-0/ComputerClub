package src

import "errors"

var EventFormatError = errors.New("event format error")
var YouShallNotPass = errors.New("YouShallNotPass")
var NotOpenYet = errors.New("NotOpenYet")
var PlaceIsBusy = errors.New("PlaceIsBusy")
var ClientUnknown = errors.New("ClientUnknown")
var ICanWaitNoLonger = errors.New("ICanWaitNoLonger!")

var QueueIsFull = errors.New("queue is full")
var QueueIsEmpty = errors.New("queue is empty")
var ClientNotInQueue = errors.New("queue does not contain client")

var DeskIsOccupied = errors.New("desk is occupied")
var DeskAlreadyFree = errors.New("desk is already free")
