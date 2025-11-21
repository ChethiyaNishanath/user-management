package restclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	rest "user-management/internal/rest-client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	BASE_URL = "https://webhook.site/e4c22ed1-5087-43c4-986e-2aa11f653189"
)

func TestRestClient_Get(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/hello")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))

	defer mockServer.Close()

	restClient := rest.NewRestClient(mockServer.URL, 1*time.Second)

	ctx := context.Background()

	requestOpts := rest.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var response any

	err := restClient.Get(ctx, "/hello", requestOpts, &response)

	if err != nil {
		panic(err)
	}
}

func TestRestClient_Post(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Path, "/hello")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))

	defer mockServer.Close()

	restClient := rest.NewRestClient(mockServer.URL, 1*time.Second)

	ctx := context.Background()

	requestOpts := rest.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"id":   "1234",
			"name": "Chethiya",
		},
	}

	var response any

	err := restClient.Post(ctx, "/hello", requestOpts, &response)

	if err != nil {
		panic(err)
	}

}

func TestRestClient_Put(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Contains(t, r.URL.Path, "/hello")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))

	defer mockServer.Close()

	restClient := rest.NewRestClient(mockServer.URL, 1*time.Second)

	ctx := context.Background()

	requestOpts := rest.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"id":   "1234",
			"name": "Chethiya",
		},
	}

	var response any

	err := restClient.Put(ctx, "/hello", requestOpts, &response)

	if err != nil {
		panic(err)
	}

}

func TestRestClient_Patch(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Contains(t, r.URL.Path, "/hello")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))

	defer mockServer.Close()

	restClient := rest.NewRestClient(mockServer.URL, 1*time.Second)

	ctx := context.Background()

	requestOpts := rest.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"id":   "1234",
			"name": "Chethiya",
		},
	}

	var response any

	err := restClient.Patch(ctx, "/hello", requestOpts, &response)

	if err != nil {
		panic(err)
	}

}

func TestRestClient_Delete(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/hello")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))

	defer mockServer.Close()

	restClient := rest.NewRestClient(mockServer.URL, 1*time.Second)

	ctx := context.Background()

	requestOpts := rest.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var response any

	err := restClient.Delete(ctx, fmt.Sprintf("/hello/%s", uuid.New()), requestOpts, &response)

	if err != nil {
		panic(err)
	}

}

func TestRestCLient_Response_Get(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/hello")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))

	defer mockServer.Close()

	restClient := rest.NewRestClient(mockServer.URL, 1*time.Second)

	ctx := context.Background()

	requestOpts := rest.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	response := make(map[string]any)

	err := restClient.Get(ctx, fmt.Sprintf("/hello/%s", uuid.New()), requestOpts, &response)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Hello", response["message"])

}
