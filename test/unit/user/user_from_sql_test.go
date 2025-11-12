package user

import (
	"testing"

	"user-management/internal/db/sqlc"
	"user-management/internal/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFromSQLC_ValidData(t *testing.T) {
	validUUID := uuid.New()

	sqlcUser := sqlc.User{
		UserID:    validUUID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "+1234567890",
		Age:       30,
		Status:    "Active",
	}

	result := user.FromSQLC(sqlcUser)

	assert.Equal(t, validUUID, result.UserId)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "john.doe@example.com", result.Email)
	assert.Equal(t, "+1234567890", result.Phone)
	assert.Equal(t, int16(30), result.Age)
	assert.Equal(t, user.Active, result.Status)
}

func TestFromSQLC_InvalidUUID(t *testing.T) {
	sqlcUser := sqlc.User{
		UserID:    uuid.Nil,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "+1234567890",
		Age:       30,
		Status:    "Active",
	}

	result := user.FromSQLC(sqlcUser)

	assert.Equal(t, uuid.Nil, result.UserId)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, user.Active, result.Status)
}

func TestFromSQLC_InvalidStatus(t *testing.T) {
	validUUID := uuid.New()

	sqlcUser := sqlc.User{
		UserID:    validUUID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "+1234567890",
		Age:       30,
		Status:    "InvalidStatus",
	}

	result := user.FromSQLC(sqlcUser)

	assert.Equal(t, validUUID, result.UserId)
	assert.Equal(t, user.InActive, result.Status)
}

func TestFromSQLC_EmptyStatus(t *testing.T) {
	validUUID := uuid.New()

	sqlcUser := sqlc.User{
		UserID:    validUUID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "+1234567890",
		Age:       30,
		Status:    "",
	}

	result := user.FromSQLC(sqlcUser)

	assert.Equal(t, user.InActive, result.Status)
}

func TestFromSQLC_InActiveStatus(t *testing.T) {
	validUUID := uuid.New()

	sqlcUser := sqlc.User{
		UserID:    validUUID,
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
		Phone:     "+9876543210",
		Age:       25,
		Status:    "InActive",
	}

	result := user.FromSQLC(sqlcUser)

	assert.Equal(t, validUUID, result.UserId)
	assert.Equal(t, user.InActive, result.Status)
}

func TestFromSQLC_SpecialCharactersInFields(t *testing.T) {
	validUUID := uuid.New()

	sqlcUser := sqlc.User{
		UserID:    validUUID,
		FirstName: "João",
		LastName:  "O'Brien-Smith",
		Email:     "test+tag@example.com",
		Phone:     "+1-234-567-8900",
		Age:       28,
		Status:    "Active",
	}

	result := user.FromSQLC(sqlcUser)

	assert.Equal(t, "João", result.FirstName)
	assert.Equal(t, "O'Brien-Smith", result.LastName)
	assert.Equal(t, "test+tag@example.com", result.Email)
	assert.Equal(t, "+1-234-567-8900", result.Phone)
}

func TestFromSQLC_CaseSensitiveStatus(t *testing.T) {
	testCases := []struct {
		name           string
		status         string
		expectedStatus user.UserStatus
	}{
		{
			name:           "Lowercase active",
			status:         "active",
			expectedStatus: user.InActive,
		},
		{
			name:           "Uppercase ACTIVE",
			status:         "ACTIVE",
			expectedStatus: user.InActive,
		},
		{
			name:           "Correct case Active",
			status:         "Active",
			expectedStatus: user.Active,
		},
		{
			name:           "Correct case InActive",
			status:         "InActive",
			expectedStatus: user.InActive,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validUUID := uuid.New()
			sqlcUser := sqlc.User{
				UserID:    validUUID,
				FirstName: "Test",
				LastName:  "User",
				Email:     "test@example.com",
				Phone:     "+1234567890",
				Age:       25,
				Status:    tc.status,
			}

			result := user.FromSQLC(sqlcUser)
			assert.Equal(t, tc.expectedStatus, result.Status)
		})
	}
}
