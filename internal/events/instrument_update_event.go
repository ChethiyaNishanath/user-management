package events

import "time"

type InstrumentUpdatedEvent struct {
	Symbol    string    `json:"symbol"`
	Price     string    `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
}
