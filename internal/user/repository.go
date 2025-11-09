package user

import (
	"context"
	"database/sql"
	"user-management/internal/db/sqlc"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) *Repository {
	return &Repository{queries: q}
}

func (r *Repository) Create(ctx context.Context, user *User) (sqlc.User, error) {

	params := sqlc.CreateUserParams{
		UserID:    user.UserId.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       user.Age,
		Status:    user.Status.String(),
	}

	return r.queries.CreateUser(ctx, params)
}

func (r *Repository) GetAll(ctx context.Context) ([]sqlc.User, error) {
	return r.queries.ListAllUsers(ctx)
}

func (r *Repository) GetUserById(ctx context.Context, userId string) (sqlc.User, error) {
	return r.queries.FindUserById(ctx, userId)
}

func (r *Repository) Update(ctx context.Context, user *User) (sqlc.User, error) {

	parms := sqlc.UpdateUserParams{
		FirstName: nullableString(user.FirstName),
		LastName:  nullableString(user.LastName),
		Email:     nullableString(user.Email),
		Phone:     nullableString(user.Phone),
		Age:       nullableInt16(user.Age),
		Status:    nullableString(user.Status.String()),
		UserID:    user.UserId.String(),
	}

	return r.queries.UpdateUser(ctx, parms)
}

func (r *Repository) Delete(ctx context.Context, userId string) error {
	return r.queries.DeleteUserByID(ctx, userId)
}

func nullableString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func nullableInt16(i int16) sql.NullInt16 {
	return sql.NullInt16{
		Int16: i,
		Valid: i > 0,
	}
}
