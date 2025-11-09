package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, u *User) (User, error) {
	newUser := NewUser(u.FirstName, u.LastName, u.Email, u.Phone, u.Age)
	savedUser, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return User{}, err
	}
	return FromSQLC(savedUser), nil
}

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return FromSQLCList(users), nil
}

func (s *Service) GetUserById(ctx context.Context, userId string) (User, error) {
	u, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		return User{}, err
	}
	return FromSQLC(u), nil
}

func (s *Service) UpdateUser(ctx context.Context, userId string, u *UserUpdateRequest) (User, error) {

	existing, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		return User{}, err
	}

	if u.FirstName != "" {
		existing.FirstName = u.FirstName
	}
	if u.LastName != "" {
		existing.LastName = u.LastName
	}
	if u.Email != "" {
		existing.Email = u.Email
	}
	if u.Phone != "" {
		existing.Phone = u.Phone
	}
	if u.Age > 0 {
		existing.Age = u.Age
	}
	if u.Status != "" {
		existing.Status = u.Status
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		return User{}, fmt.Errorf("invalid userId: %w", err)
	}

	userToBeUpdate, err := UpdatableUser(id, &existing)

	if err != nil {
		return User{}, err
	}

	savedUser, err := s.repo.Update(ctx, userToBeUpdate)
	if err != nil {
		return User{}, err
	}
	return FromSQLC(savedUser), nil
}

func (s *Service) DeleteUserById(ctx context.Context, userId string) error {
	err := s.repo.Delete(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}
