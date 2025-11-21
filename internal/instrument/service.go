package instrument

import (
	"context"
	"fmt"
	"user-management/internal/common/converters"
	events "user-management/internal/events"
	ws "user-management/internal/ws"

	"github.com/google/uuid"
)

type Service struct {
	repo     *Repository
	eventBus *events.Bus
}

func NewService(repo *Repository, bus *events.Bus) *Service {
	return &Service{repo: repo, eventBus: bus}
}

func (s *Service) CreateInstrument(ctx context.Context, i *Instrument) (Instrument, error) {
	newInstrument := NewInstrument(i.Symbol, i.Name, i.InstrumentType, i.Exchange, i.LastPrice)
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

func (s *Service) GetInstrumentBySymbol(ctx context.Context, symbol string) (Instrument, error) {
	instrument, err := s.repo.GetInstrumentBySymbol(ctx, symbol)
	if err != nil {
		return Instrument{}, err
	}
	return FromSQLC(instrument), nil
}

func (s *Service) UpdateInstrument(ctx context.Context, instrumentId string, i *InstrumentUpdateRequest, connMgr *ws.ConnectionManager) (Instrument, error) {

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
	if i.InstrumentType != "" {
		existing.InstrumentType = i.InstrumentType
	}
	if i.Exchange != "" {
		existing.Exchange = i.Exchange
	}

	priceChanged := false
	if i.LastPrice > 0 {
		newPriceStr := converters.Float64ToString(i.LastPrice)
		if existing.LastPrice != newPriceStr {
			priceChanged = true
			existing.LastPrice = newPriceStr
		}
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

	if priceChanged {

		evt := events.InstrumentUpdatedEvent{
			Symbol:    savedInstrument.Symbol,
			Price:     savedInstrument.LastPrice,
			UpdatedAt: savedInstrument.UpdatedAt,
		}

		s.eventBus.Publish(events.InstrumentUpdated, evt)
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
