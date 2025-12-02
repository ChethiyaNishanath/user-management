package instrument

import (
	"context"
	"log/slog"
	"user-management/internal/common/converters"
	"user-management/internal/db/sqlc"

	"github.com/google/uuid"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) *Repository {
	return &Repository{queries: q}
}

func (r *Repository) Create(ctx context.Context, instrument *Instrument) (sqlc.Instrument, error) {

	params := sqlc.CreateInstrumentParams{
		Symbol:         instrument.Symbol,
		Name:           instrument.Name,
		InstrumentType: instrument.Instrument_Type,
		Exchange:       instrument.Exchange,
		LastPrice:      converters.Float64ToString(instrument.Last_Price),
		UpdatedAt:      instrument.Updated_At,
		ID:             instrument.Id,
	}

	return r.queries.CreateInstrument(ctx, params)
}

func (r *Repository) GetAllPaged(ctx context.Context, limit int, offset int) ([]sqlc.Instrument, error) {

	params := sqlc.ListAllInstrumentPagedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	return r.queries.ListAllInstrumentPaged(ctx, params)
}

func (r *Repository) GetInstrumentById(ctx context.Context, instrumentId string) (sqlc.Instrument, error) {

	parsedUUID, err := uuid.Parse(instrumentId)
	if err != nil {
		slog.Error("Invalid UUID from DB", "error", err)
		return sqlc.Instrument{}, err
	}

	return r.queries.FindInstrumentById(ctx, parsedUUID)
}

func (r *Repository) Update(ctx context.Context, instrument *Instrument) (sqlc.Instrument, error) {

	parms := sqlc.UpdateInstrumentParams{
		Symbol:         converters.NullableString(instrument.Symbol),
		Name:           converters.NullableString(instrument.Name),
		InstrumentType: converters.NullableString(instrument.Instrument_Type),
		Exchange:       converters.NullableString(instrument.Exchange),
		LastPrice:      converters.NullableFloat64(instrument.Last_Price),
		UpdatedAt:      converters.NullableTime(instrument.Updated_At),
		ID:             instrument.Id,
	}

	return r.queries.UpdateInstrument(ctx, parms)
}

func (r *Repository) Delete(ctx context.Context, instrumentId string) error {

	parsedUUID, err := uuid.Parse(instrumentId)
	if err != nil {
		slog.Error("Invalid UUID from DB", "error", err)
		return err
	}

	return r.queries.DeleteInstrumentById(ctx, parsedUUID)
}
