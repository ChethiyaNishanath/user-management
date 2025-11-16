package app

import (
	"context"
	"database/sql"
	"user-management/internal/db/sqlc"
	events "user-management/internal/events"
	"user-management/internal/instrument"
	"user-management/internal/middleware"
	"user-management/internal/user"
	"user-management/internal/validation"
	ws "user-management/internal/ws"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	DB        *sql.DB
	Validator *validator.Validate
	Queries   *sqlc.Queries

	UserHandler       *user.Handler
	InstrumentHandler *instrument.Handler
	WebSocketHandler  *ws.Handler
}

func NewApp(db *sql.DB, ctx *context.Context) *App {
	validate := validator.New()
	validation.RegisterValidations(validate)

	queries := sqlc.New(db)
	bus := events.NewBus()
	connMgr := ws.NewConnectionManager(*ctx)

	websocketModule := ws.NewModule(ctx, connMgr, bus)
	instrumentModule := instrument.NewModule(queries, validate, connMgr, bus)
	userModule := user.NewModule(queries, validate)

	return &App{
		DB:                db,
		Validator:         validate,
		Queries:           queries,
		UserHandler:       userModule.Handler,
		InstrumentHandler: instrumentModule.Handler,
		WebSocketHandler:  websocketModule.Handler,
	}
}

func (a *App) RegisterRoutes(r chi.Router) {

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", a.UserHandler.CreateUser)
		r.With(middleware.Paginate).Get("/", a.UserHandler.GetUsers)
		r.Get("/{id}", a.UserHandler.GetUserById)
		r.Patch("/{id}", a.UserHandler.UpdateUserById)
		r.Delete("/{id}", a.UserHandler.DeleteUserById)
	})

	r.Route("/instruments", func(r chi.Router) {
		r.Post("/", a.InstrumentHandler.CreateInstrument)
		r.With(middleware.Paginate).Get("/", a.InstrumentHandler.GetInstruments)
		r.Get("/{id}", a.InstrumentHandler.GetInstrumentById)
		r.Patch("/{id}", a.InstrumentHandler.UpdateInstrumentById)
		r.Delete("/{id}", a.InstrumentHandler.DeleteInstrumentById)
	})

	r.Get("/ws", a.WebSocketHandler.HandleWebSocket)
}
