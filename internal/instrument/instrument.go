package instrument

import (
	"log/slog"
	"strconv"
	"time"
	"user-management/internal/common/converters"
	"user-management/internal/db/sqlc"

	"github.com/google/uuid"
)

type Instrument struct {
	Id              uuid.UUID `json:"id"`
	Symbol          string    `json:"symbol" validate:"required,min=2,max=50"`
	Name            string    `json:"name" validate:"required,min=2,max=50"`
	Instrument_Type string    `json:"type" validate:"required,email"`
	Exchange        string    `json:"exchange" validate:"omitempty,e164"`
	Last_Price      float64   `json:"last_price" validate:"omitempty,gt=0"`
	Created_At      time.Time `json:"created_At" validate:"omitempty,userStatus"`
	Updated_At      time.Time `json:"updated_At" validate:"omitempty,userStatus"`
}

func NewInstrument(symbol string, name string, instrumentType string, exchange string, lastPrice float64) *Instrument {
	return &Instrument{
		Id:              uuid.New(),
		Symbol:          symbol,
		Name:            name,
		Instrument_Type: instrumentType,
		Exchange:        exchange,
		Last_Price:      lastPrice,
		Created_At:      time.Now(),
		Updated_At:      time.Now(),
	}
}

func UpdatableInstrument(id uuid.UUID, instrument *sqlc.Instrument) (*Instrument, error) {

	lastPrice, err := converters.StringToFloat64(instrument.LastPrice)

	if err != nil {
		return nil, err
	}

	return &Instrument{
		Id:              id,
		Symbol:          instrument.Symbol,
		Name:            instrument.Name,
		Instrument_Type: instrument.InstrumentType,
		Exchange:        instrument.Exchange,
		Last_Price:      lastPrice,
		Created_At:      instrument.CreatedAt,
		Updated_At:      time.Now(),
	}, nil
}

func FromSQLC(i sqlc.Instrument) Instrument {

	lastPrice, err := strconv.ParseFloat(i.LastPrice, 64)

	if err != nil {
		slog.Error("Error parsing string to float64", "error", err)
	}

	return Instrument{
		Id:              i.ID,
		Symbol:          i.Symbol,
		Name:            i.Name,
		Instrument_Type: i.InstrumentType,
		Exchange:        i.Exchange,
		Last_Price:      float64(lastPrice),
		Created_At:      i.CreatedAt,
		Updated_At:      i.UpdatedAt,
	}
}

func FromSQLCList(instruments []sqlc.Instrument) []Instrument {
	mapped := make([]Instrument, len(instruments))
	for i, u := range instruments {
		mapped[i] = FromSQLC(u)
	}
	return mapped
}
