package event_manager

type ClientEvent struct {
	Time       string `json:"time"`
	ClientName string `json:"client_name"`
	EventID    string `json:"event_id"`
}

type ClientArrival struct {
	Time       string `json:"time"`
	ClientName string `json:"client_name"`
}

type ClientDeparture struct {
	Time       string `json:"time"`
	ClientName string `json:"client_name"`
}

type TableOccupation struct {
	Time       string `json:"time"`
	ClientName string `json:"client_name"`
	TableID    string `json:"table_id"`
}

type ClientAwait struct {
	Time       string `json:"time"`
	ClientName string `json:"client_name"`
}
