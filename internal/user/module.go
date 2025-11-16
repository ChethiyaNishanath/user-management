package user

import (
	"user-management/internal/db/sqlc"

	"github.com/go-playground/validator/v10"
)

type Module struct {
	Handler *Handler
	Service *Service
	Repo    *Repository
}

func NewModule(q *sqlc.Queries, v *validator.Validate) *Module {
	repo := NewRepository(q)
	service := NewService(repo)
	handler := NewHandler(service, v)

	return &Module{
		Handler: handler,
		Service: service,
		Repo:    repo,
	}
}
