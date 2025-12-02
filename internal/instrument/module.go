// package instrument

// import (
// 	"context"
// 	"user-management/internal/events"
// 	"user-management/internal/ws"
// )

// func RegisterModule(
// 	ctx *context.Context,
// 	service *Service,
// 	wsRouter *ws.Router,
// 	connMgr *ws.ConnectionManager,
// 	eventBus *events.Bus,
// ) {
// 	RegisterWsRoutes(wsRouter, service)

// 	eventBus.Subscribe(events.InstrumentUpdated, func(e events.Event) {
// 		evt := e.Data.(events.InstrumentUpdatedEvent)
// 		connMgr.Broadcast(*ctx, evt.Symbol, evt)
// 	})
// }

package instrument

import (
	"context"
	wsutils "user-management/internal/common/wsutils"
	"user-management/internal/db/sqlc"
	events "user-management/internal/events"
	ws "user-management/internal/ws"

	"github.com/go-playground/validator/v10"
)

type Module struct {
	Handler *Handler
	Service *Service
	Repo    *Repository
}

func NewModule(q *sqlc.Queries, v *validator.Validate, connMger *ws.ConnectionManager, bus *events.Bus) *Module {
	repo := NewRepository(q)
	service := NewService(repo, bus)
	handler := NewHandler(service, v, connMger)

	module := &Module{
		Handler: handler,
		Service: service,
		Repo:    repo,
	}

	module.registerEventSubscribers(bus)
	return module
}

func (m *Module) registerEventSubscribers(bus *events.Bus) {
	bus.Subscribe(events.InstrumentUpdated, func(e events.Event) {
		evt := e.Data.(events.InstrumentUpdatedEvent)

		m.Handler.ConnMgr.Broadcast(
			context.Background(),
			evt.Symbol,
			wsutils.WSMessage{
				Action: e.Action,
				Topic:  e.Topic,
				Data:   evt,
			},
		)
	})
}
