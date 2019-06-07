package contracts

import "time"

type EventCreatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (e *EventCreatedEvent) EventName() string {
	return "event.created"
}
