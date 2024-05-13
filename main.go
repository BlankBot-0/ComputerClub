package main

import (
	"yadroTest/internal/event_manager"
	"yadroTest/internal/impl"
)

func main() {
	cp := impl.NewClientPool()
	cq := impl.NewClientListQueue()
	deskManager := impl.NewDesks(5)

	event_manager.NewEventManager(event_manager.Deps{
		ClientPool:  cp,
		DeskManager: deskManager,
		ClientQueue: cq,
	})

}
