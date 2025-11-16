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
	Id             uuid.UUID `json:"id"`
	Symbol         string    `json:"symbol" validate:"required,min=1,max=20"`
	Name           string    `json:"name" validate:"required,min=2,max=100"`
	InstrumentType string    `json:"instrument_type" validate:"required,oneof=Equity Bond ETF Crypto Forex"`
	Exchange       string    `json:"exchange" validate:"required,alphanum,min=2,max=10"`
	LastPrice      float64   `json:"last_price" validate:"required,gt=0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewInstrument(symbol string, name string, instrumentType string, exchange string, lastPrice float64) *Instrument {
	return &Instrument{
		Id:             uuid.New(),
		Symbol:         symbol,
		Name:           name,
		InstrumentType: instrumentType,
		Exchange:       exchange,
		LastPrice:      lastPrice,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func UpdatableInstrument(id uuid.UUID, instrument *sqlc.Instrument) (*Instrument, error) {

	lastPrice, err := converters.StringToFloat64(instrument.LastPrice)

	if err != nil {
		return nil, err
	}

	return &Instrument{
		Id:             id,
		Symbol:         instrument.Symbol,
		Name:           instrument.Name,
		InstrumentType: instrument.InstrumentType,
		Exchange:       instrument.Exchange,
		LastPrice:      lastPrice,
		CreatedAt:      instrument.CreatedAt,
		UpdatedAt:      time.Now(),
	}, nil
}

func FromSQLC(i sqlc.Instrument) Instrument {

	lastPrice, err := strconv.ParseFloat(i.LastPrice, 64)

	if err != nil {
		slog.Error("Error parsing string to float64", "error", err)
	}

	return Instrument{
		Id:             i.ID,
		Symbol:         i.Symbol,
		Name:           i.Name,
		InstrumentType: i.InstrumentType,
		Exchange:       i.Exchange,
		LastPrice:      float64(lastPrice),
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}

func FromSQLCList(instruments []sqlc.Instrument) []Instrument {
	mapped := make([]Instrument, len(instruments))
	for i, u := range instruments {
		mapped[i] = FromSQLC(u)
	}
	return mapped
}
