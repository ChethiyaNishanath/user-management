package user_test

import (
	"testing"

	"user-management/internal/user"

	"github.com/stretchr/testify/assert"
)

func TestUserStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status user.UserStatus
		want   string
	}{
		{
			name:   "Active status returns 'Active'",
			status: user.Active,
			want:   "Active",
		},
		{
			name:   "InActive status returns 'InActive'",
			status: user.InActive,
			want:   "InActive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseUserStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    user.UserStatus
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Parse 'Active' string",
			input:   "Active",
			want:    user.Active,
			wantErr: false,
		},
		{
			name:    "Parse 'InActive' string",
			input:   "InActive",
			want:    user.InActive,
			wantErr: false,
		},
		{
			name:    "Invalid status string",
			input:   "Invalid",
			want:    0,
			wantErr: true,
			errMsg:  "invalid user status: Invalid",
		},
		{
			name:    "Empty string",
			input:   "",
			want:    0,
			wantErr: true,
			errMsg:  "invalid user status: ",
		},
		{
			name:    "Lowercase 'active'",
			input:   "active",
			want:    0,
			wantErr: true,
			errMsg:  "invalid user status: active",
		},
		{
			name:    "Uppercase 'ACTIVE'",
			input:   "ACTIVE",
			want:    0,
			wantErr: true,
			errMsg:  "invalid user status: ACTIVE",
		},
		{
			name:    "Lowercase 'inactive'",
			input:   "inactive",
			want:    0,
			wantErr: true,
			errMsg:  "invalid user status: inactive",
		},
		{
			name:    "Random string",
			input:   "random",
			want:    0,
			wantErr: true,
			errMsg:  "invalid user status: random",
		},
		{
			name:    "Numeric string",
			input:   "123",
			want:    0,
			wantErr: true,
			errMsg:  "invalid user status: 123",
		},
		{
			name:    "String with spaces",
			input:   " Active ",
			want:    0,
			wantErr: true,
			errMsg:  "invalid user status:  Active ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.ParseUserStatus(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
				if tt.errMsg != "" {
					assert.EqualError(t, err, tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseUserStatus_ValidCases(t *testing.T) {
	validCases := map[string]user.UserStatus{
		"Active":   user.Active,
		"InActive": user.InActive,
	}

	for input, expected := range validCases {
		t.Run("Valid_"+input, func(t *testing.T) {
			got, err := user.ParseUserStatus(input)
			assert.NoError(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestParseUserStatus_InvalidCases(t *testing.T) {
	invalidInputs := []string{
		"",
		"invalid",
		"ACTIVE",
		"active",
		"INACTIVE",
		"inactive",
		"pending",
		"suspended",
		"0",
		"1",
		" Active",
		"Active ",
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input, func(t *testing.T) {
			got, err := user.ParseUserStatus(input)
			assert.Error(t, err)
			assert.Equal(t, user.UserStatus(0), got)
			assert.Contains(t, err.Error(), "invalid user status")
		})
	}
}

func TestParseUserStatus_ErrorMessage(t *testing.T) {
	input := "UnknownStatus"
	_, err := user.ParseUserStatus(input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user status")
	assert.Contains(t, err.Error(), input)
}

func TestParseUserStatus_RoundTrip(t *testing.T) {
	statuses := []user.UserStatus{user.Active, user.InActive}

	for _, status := range statuses {
		t.Run("RoundTrip_"+status.String(), func(t *testing.T) {
			str := status.String()
			parsed, err := user.ParseUserStatus(str)
			assert.NoError(t, err)
			assert.Equal(t, status, parsed)
		})
	}
}
