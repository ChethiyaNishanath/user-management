package helper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewJSONRequest(t *testing.T, method, url string, body any) *http.Request {
	b, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(method, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return req
}
