package stringutils

import (
	"crypto/subtle"
	"strings"
)

func Equals(a, b string) bool {
	return a == b
}

func EqualsIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}

func SafeEquals(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func IsEmpty(s string) bool {
	return len(s) == 0
}

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func ContainsIgnoreCase(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}
