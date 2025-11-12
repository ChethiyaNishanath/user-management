package user

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

func (r *Repository) Create(ctx context.Context, user *User) (sqlc.User, error) {

	params := sqlc.CreateUserParams{
		UserID:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       user.Age,
		Status:    user.Status.String(),
	}

	return r.queries.CreateUser(ctx, params)
}

func (r *Repository) GetAllPaged(ctx context.Context, limit int, offset int) ([]sqlc.User, error) {

	params := sqlc.ListAllUsersPagedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	return r.queries.ListAllUsersPaged(ctx, params)
}

func (r *Repository) GetUserById(ctx context.Context, userId string) (sqlc.User, error) {

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("Invalid UUID from DB", "error", err)
		return sqlc.User{}, err
	}

	return r.queries.FindUserById(ctx, parsedUUID)
}

func (r *Repository) Update(ctx context.Context, user *User) (sqlc.User, error) {

	parms := sqlc.UpdateUserParams{
		FirstName: converters.NullableString(user.FirstName),
		LastName:  converters.NullableString(user.LastName),
		Email:     converters.NullableString(user.Email),
		Phone:     converters.NullableString(user.Phone),
		Age:       converters.NullableInt16(user.Age),
		Status:    converters.NullableString(user.Status.String()),
		UserID:    user.UserId,
	}

	return r.queries.UpdateUser(ctx, parms)
}

func (r *Repository) Delete(ctx context.Context, userId string) error {

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("Invalid UUID from DB", "error", err)
		return err
	}

	return r.queries.DeleteUserByID(ctx, parsedUUID)
}
