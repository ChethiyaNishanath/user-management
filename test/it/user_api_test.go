package it

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"user-management/internal/app"
	"user-management/internal/db"
	"user-management/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var r *chi.Mux

func TestMain(m *testing.M) {

	ctx := context.Background()

	dbName := "usermanagementdb"
	dbUser := "postgres"
	dbPassword := "Test12344"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)

	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			slog.Error(fmt.Sprintf("failed to terminate container: %s", err))
		}
	}()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to start container: %s", err))
		return
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		slog.Error("failed to get connection string: %v", "error", err)
	}

	dbConn := db.Connect(connStr)
	defer dbConn.Close()

	newApp := app.NewApp(dbConn)

	r := chi.NewRouter()
	newApp.RegisterRoutes(r)
}

func TestCreateUserAPI(t *testing.T) {
	reqBody := `{
		"firstName": "Chethiya",
        "lastName": "Viharagama",
        "email": "chethiya.viharagama@yaalalabs.com",
        "phone": "+94768680618",
        "age": 11,
        "status": "Active"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUser user.User
	err := json.NewDecoder(w.Body).Decode(&createdUser)
	require.NoError(t, err, "Failed to decode response")

	assert.NotEqual(t, uuid.Nil, createdUser.UserId, "User ID should not be nil")
	assert.Equal(t, "Chethiya", createdUser.FirstName)
	assert.Equal(t, "Viharagama", createdUser.LastName)
	assert.Equal(t, "chethiya.viharagama@yaalalabs.com", createdUser.Email)
	assert.Equal(t, "+94768680618", createdUser.Phone)
	assert.Equal(t, int16(11), createdUser.Age)
	assert.Equal(t, user.Active, createdUser.Status)
}

func TestGetAllUsersAPI(t *testing.T) {
	reqBody := `{
		"firstName": "Test",
        "lastName": "User",
        "email": "test.user@example.com",
        "phone": "+94768680619",
        "age": 25,
        "status": "Active"
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

	var users []user.User
	err := json.NewDecoder(w.Body).Decode(&users)
	require.NoError(t, err, "Failed to decode response")

	assert.Greater(t, len(users), 0, "Should have at least one user")
}

func TestGetUserByIdAPI(t *testing.T) {
	reqBody := `{
		"firstName": "John",
        "lastName": "Doe",
        "email": "john.doe@example.com",
        "phone": "+94768680620",
        "age": 30,
        "status": "Active"
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	var createdUser user.User
	err := json.NewDecoder(createW.Body).Decode(&createdUser)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", createdUser.UserId.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

	var fetchedUser user.User
	err = json.NewDecoder(w.Body).Decode(&fetchedUser)
	require.NoError(t, err, "Failed to decode response")

	assert.Equal(t, createdUser.UserId, fetchedUser.UserId)
	assert.Equal(t, "John", fetchedUser.FirstName)
	assert.Equal(t, "Doe", fetchedUser.LastName)
	assert.Equal(t, "john.doe@example.com", fetchedUser.Email)
	assert.Equal(t, "+94768680620", fetchedUser.Phone)
	assert.Equal(t, int16(30), fetchedUser.Age)
	assert.Equal(t, user.Active, fetchedUser.Status)
}

func TestGetUserByIdAPI_NotFound(t *testing.T) {
	nonExistentID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", nonExistentID.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status 404 Not Found")
}

func TestGetUserByIdAPI_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 Bad Request")
}

func TestUpdateUserAPI(t *testing.T) {
	reqBody := `{
		"firstName": "Jhon",
        "lastName": "Doe",
        "email": "jane.doe@example.com",
        "phone": "+94768680622",
        "age": 29,
        "status": "InActive"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUser user.User
	err := json.NewDecoder(w.Body).Decode(&createdUser)
	require.NoError(t, err, "Failed to decode response")

	updateReqBody := `{
		"firstName": "Jane",
        "lastName": "Doe",
        "email": "jane.doe@example.com",
        "phone": "+94768680622",
        "age": 29,
        "status": "InActive"
	}`

	updateReq := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%s", createdUser.UserId.String()), strings.NewReader(updateReqBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)

	assert.Equal(t, http.StatusOK, updateW.Code, "Expected status 200 OK")

	var updatedUser user.User
	err = json.NewDecoder(updateW.Body).Decode(&updatedUser)
	require.NoError(t, err)

	assert.Equal(t, createdUser.UserId, updatedUser.UserId)
	assert.Equal(t, "Jane", updatedUser.FirstName)
	assert.Equal(t, "Doe", updatedUser.LastName)
	assert.Equal(t, "jane.doe@example.com", updatedUser.Email)
	assert.Equal(t, "+94768680622", updatedUser.Phone)
	assert.Equal(t, int16(29), updatedUser.Age)
	assert.Equal(t, user.InActive, updatedUser.Status)

}

func TestDeleteUserAPI(t *testing.T) {
	reqBody := `{
		"firstName": "Jhon",
        "lastName": "Doe",
        "email": "jane.doe@example.com",
        "phone": "+94768680622",
        "age": 29,
        "status": "InActive"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUser user.User
	err := json.NewDecoder(w.Body).Decode(&createdUser)
	require.NoError(t, err, "Failed to decode response")

	updateReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", createdUser.UserId.String()), nil)
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)
}

func TestDeleteUserWithMultipleUsersAPI(t *testing.T) {
	reqBodyOne := `{
		"firstName": "Jhon",
        "lastName": "Doe",
        "email": "jane.doe@example.com",
        "phone": "+94768680622",
        "age": 29,
        "status": "InActive"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBodyOne))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUserOne user.User
	err := json.NewDecoder(w.Body).Decode(&createdUserOne)
	require.NoError(t, err, "Failed to decode response")

	//========================================================================

	reqBodyTwo := `{
		"firstName": "Chethiya",
        "lastName": "Nishanath",
        "email": "chethiya.nishanath@example.com",
        "phone": "+94768680600",
        "age": 29,
        "status": "InActive"
	}`

	reqTwo := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBodyTwo))
	reqTwo.Header.Set("Content-Type", "application/json")
	wTwo := httptest.NewRecorder()

	r.ServeHTTP(wTwo, reqTwo)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUserTwo user.User
	errTwo := json.NewDecoder(wTwo.Body).Decode(&createdUserTwo)
	require.NoError(t, errTwo, "Failed to decode response")

	updateReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", createdUserOne.UserId.String()), nil)
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)

	//=========================================================

	reqAll := httptest.NewRequest(http.MethodGet, "/users", nil)
	wAll := httptest.NewRecorder()
	r.ServeHTTP(wAll, reqAll)

	assert.Equal(t, http.StatusOK, wAll.Code, "Expected status 200 OK")

	var users []user.User
	errAll := json.NewDecoder(w.Body).Decode(&users)

	require.NoError(t, errAll, "Failed to decode response")

	assert.Greater(t, len(users), 0, "Should have at least one user")
}
