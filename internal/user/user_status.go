package user

import (
	"encoding/json"
	"fmt"
)

type UserStatus int

const (
	Active UserStatus = iota
	InActive
)

var stateName = map[UserStatus]string{
	Active:   "Active",
	InActive: "InActive",
}

func (s UserStatus) String() string {
	return stateName[s]
}

func ParseUserStatus(s string) (UserStatus, error) {
	switch s {
	case "Active":
		return Active, nil
	case "InActive":
		return InActive, nil
	}
	return 0, fmt.Errorf("invalid user status: %s", s)
}

func (s UserStatus) MarshalJSON() ([]byte, error) {
	switch s {
	case 0:
		return json.Marshal("Active")
	case 1:
		return json.Marshal("Inactive")
	default:
		return nil, fmt.Errorf("unknown status value")
	}
}

func (s *UserStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case "Active":
		*s = 0
	case "Inactive":
		*s = 1
	default:
		return fmt.Errorf("invalid status string: %s", str)
	}

	return nil
}
