package app

import (
	"database/sql"
	"user-management/internal/db/sqlc"
	"user-management/internal/middleware"
	"user-management/internal/user"
	"user-management/internal/validation"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	DB        *sql.DB
	Validator *validator.Validate
	Queries   *sqlc.Queries

	UserHandler *user.Handler
}

func NewApp(db *sql.DB) *App {
	validate := validator.New()
	validation.RegisterValidations(validate)

	queries := sqlc.New(db)

	userRepo := user.NewRepository(queries)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, validate)

	return &App{
		DB:          db,
		Validator:   validate,
		Queries:     queries,
		UserHandler: userHandler,
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
}
