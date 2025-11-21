package binance

import (
	"context"
	"strings"
	"user-management/internal/config"
	"user-management/internal/events"
	ws "user-management/internal/ws"
)

type Module struct {
	ConnMgr *ws.ConnectionManager
}

func NewModule(ctx *context.Context, bus *events.Bus, connMger *ws.ConnectionManager) *Module {

	module := &Module{
		ConnMgr: connMger,
	}

	module.registerEventSubscribers(bus)
	return module

}

func (m *Module) registerEventSubscribers(bus *events.Bus) {

	config := config.GetConfig()

	symbols := strings.Split(strings.ToLower(config.Binance.SubscribedSymbols), ",")

	binanceDepthTopics := []string{}

	for _, symbol := range symbols {
		binanceDepthTopics = append(binanceDepthTopics, symbol+"@depth")
	}

	for _, sbdt := range binanceDepthTopics {
		bus.Subscribe(sbdt, func(e events.Event) {
			evt := e.Data.(DepthStreamEvent)
			m.ConnMgr.Broadcast(
				context.Background(),
				e.Topic,
				ws.WSMessage{
					Data: evt,
				},
			)
		})
	}
}
