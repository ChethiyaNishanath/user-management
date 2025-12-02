package instrument

import (
	"context"
	"fmt"
	"user-management/internal/common/converters"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateInstrument(ctx context.Context, i *Instrument) (Instrument, error) {
	newInstrument := NewInstrument(i.Symbol, i.Name, i.Instrument_Type, i.Exchange, i.Last_Price)
	savedInstrument, err := s.repo.Create(ctx, newInstrument)
	if err != nil {
		return Instrument{}, err
	}
	return FromSQLC(savedInstrument), nil
}

func (s *Service) ListInstrumentsPaged(ctx context.Context, limit int, offset int) ([]Instrument, error) {
	instruments, err := s.repo.GetAllPaged(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return FromSQLCList(instruments), nil
}

func (s *Service) GetInstrumentById(ctx context.Context, instrumentId string) (Instrument, error) {
	u, err := s.repo.GetInstrumentById(ctx, instrumentId)
	if err != nil {
		return Instrument{}, err
	}
	return FromSQLC(u), nil
}

func (s *Service) UpdateInstrument(ctx context.Context, instrumentId string, i *InstrumentUpdateRequest) (Instrument, error) {

	existing, err := s.repo.GetInstrumentById(ctx, instrumentId)
	if err != nil {
		return Instrument{}, err
	}

	if i.Symbol != "" {
		existing.Symbol = i.Symbol
	}
	if i.Name != "" {
		existing.Name = i.Name
	}
	if i.Instrument_Type != "" {
		existing.InstrumentType = i.Instrument_Type
	}
	if i.Exchange != "" {
		existing.Exchange = i.Exchange
	}
	if i.Last_Price > 0 {
		existing.LastPrice = converters.Float64ToString(i.Last_Price)
	}

	id, err := uuid.Parse(instrumentId)
	if err != nil {
		return Instrument{}, fmt.Errorf("invalid instrumentId: %w", err)
	}

	instrumentToBeUpdate, err := UpdatableInstrument(id, &existing)

	if err != nil {
		return Instrument{}, err
	}

	savedInstrument, err := s.repo.Update(ctx, instrumentToBeUpdate)
	if err != nil {
		return Instrument{}, err
	}
	return FromSQLC(savedInstrument), nil
}

func (s *Service) DeleteInstrumentById(ctx context.Context, instrumentId string) error {
	err := s.repo.Delete(ctx, instrumentId)
	if err != nil {
		return err
	}
	return nil
}
