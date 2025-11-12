package converters

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NullableString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func NullableInt16(i int16) sql.NullInt16 {
	return sql.NullInt16{
		Int16: i,
		Valid: i > 0,
	}
}

func NullableFloat64(f float64) sql.NullString {
	if f == 0 {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: fmt.Sprintf("%f", f),
		Valid:  true,
	}
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func StringToFloat64(s string) (float64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseFloat(s, 64)
}

func NullableTime(s time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  s,
		Valid: !s.IsZero(),
	}
}
