package user

import (
	"log/slog"
	"user-management/internal/db/sqlc"

	"github.com/google/uuid"
)

type User struct {
	UserId    uuid.UUID  `json:"userId"`
	FirstName string     `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string     `json:"lastName" validate:"required,min=2,max=50"`
	Email     string     `json:"email" validate:"required,email"`
	Phone     string     `json:"phone" validate:"omitempty,e164"`
	Age       int16      `json:"age" validate:"omitempty,gt=0"`
	Status    UserStatus `json:"status" validate:"omitempty,userStatus"`
}

func NewUser(firstName string, lastName string, email string, phone string, age int16) *User {
	return &User{
		UserId:    uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Age:       age,
		Status:    Active,
	}
}

func UpdatableUser(userId uuid.UUID, u *sqlc.User) (*User, error) {
	status, err := ParseUserStatus(u.Status)

	if err != nil {
		return nil, err
	}

	return &User{
		UserId:    userId,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Age:       u.Age,
		Status:    status,
	}, nil
}

func FromSQLC(u sqlc.User) User {

	parsedUUID, err := uuid.Parse(u.UserID)
	if err != nil {
		slog.Error("Invalid UUID from DB", "error", err)
		parsedUUID = uuid.Nil
	}

	parsedStatus, err := ParseUserStatus(u.Status)
	if err != nil {
		slog.Error("Invalid status from DB, defaulting to INACTIVE", "error", u.Status)
		parsedStatus = InActive
	}

	return User{
		UserId:    parsedUUID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Age:       u.Age,
		Status:    parsedStatus,
	}
}

func FromSQLCList(users []sqlc.User) []User {
	mapped := make([]User, len(users))
	for i, u := range users {
		mapped[i] = FromSQLC(u)
	}
	return mapped
}
