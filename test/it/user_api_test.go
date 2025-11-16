package it

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"user-management/internal/user"
	helper "user-management/test/helper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {

	reqBody := map[string]interface{}{
		"firstName": "Chethiya",
		"lastName":  "Viharagama",
		"email":     "chethiya.viharagama@yaalalabs.com",
		"phone":     "+94768680618",
		"age":       30,
		"status":    "InActive",
	}

	req := helper.NewJSONRequest(t, http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUser user.User
	err := json.NewDecoder(w.Body).Decode(&createdUser)
	require.NoError(t, err, "Failed to decode response")

	assert.NotEqual(t, uuid.Nil, createdUser.UserId, "User ID should not be nil")
	assert.Equal(t, "Chethiya", createdUser.FirstName)
	assert.Equal(t, "Viharagama", createdUser.LastName)
	assert.Equal(t, "chethiya.viharagama@yaalalabs.com", createdUser.Email)
	assert.Equal(t, "+94768680618", createdUser.Phone)
	assert.Equal(t, int16(30), createdUser.Age)
	assert.Equal(t, user.Active, createdUser.Status)
}

func TestGetAllUsersAPI(t *testing.T) {
	reqBody := `{
		"firstName": "Test",
        "lastName": "User",
        "email": "test.user@example.com",
        "phone": "+94768680619",
        "age": 25
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	TestRouter.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

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
        "age": 30
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	TestRouter.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	var createdUser user.User
	err := json.NewDecoder(createW.Body).Decode(&createdUser)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", createdUser.UserId.String()), nil)
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

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
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status 404 Not Found")
}

func TestGetUserByIdAPI_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 Bad Request")
}

func TestUpdateUserAPI(t *testing.T) {

	reqBody := map[string]interface{}{
		"firstName": "Jhon",
		"lastName":  "Doe",
		"email":     "jane.doe@example.com",
		"phone":     "+94768680622",
		"age":       29,
		"status":    "InActive",
	}

	req := helper.NewJSONRequest(t, http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUser user.User
	err := json.NewDecoder(w.Body).Decode(&createdUser)
	require.NoError(t, err, "Failed to decode response")

	updateReqBody := map[string]interface{}{
		"firstName": "Jane",
		"lastName":  "Doe",
		"email":     "jane.doe@example.com",
		"phone":     "+94768680622",
		"age":       29,
		"status":    "InActive",
	}

	updateReq := helper.NewJSONRequest(t, http.MethodPatch, fmt.Sprintf("/users/%s", createdUser.UserId.String()), updateReqBody)
	updateW := httptest.NewRecorder()
	TestRouter.ServeHTTP(updateW, updateReq)

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
	reqBody := map[string]interface{}{
		"firstName": "Jhon2",
		"lastName":  "Doe2",
		"email":     "jane.doe2@example.com",
		"phone":     "+94768680622",
		"age":       29,
		"status":    "InActive",
	}

	req := helper.NewJSONRequest(t, http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUser user.User
	err := json.NewDecoder(w.Body).Decode(&createdUser)
	require.NoError(t, err, "Failed to decode response")

	updateReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", createdUser.UserId.String()), nil)
	updateW := httptest.NewRecorder()
	TestRouter.ServeHTTP(updateW, updateReq)
}

func TestDeleteUserWithMultipleUsersAPI(t *testing.T) {
	reqBodyOne := map[string]interface{}{
		"firstName": "Jhon11",
		"lastName":  "Doe11",
		"email":     "jane.doe11@example.com",
		"phone":     "+94768680622",
		"age":       29,
		"status":    "InActive",
	}

	req := helper.NewJSONRequest(t, http.MethodPost, "/users", reqBodyOne)
	w := httptest.NewRecorder()

	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUserOne user.User
	err := json.NewDecoder(w.Body).Decode(&createdUserOne)
	require.NoError(t, err, "Failed to decode response")

	//========================================================================

	reqBodyTwo := map[string]interface{}{
		"firstName": "Chethiya1",
		"lastName":  "Viharagama1",
		"email":     "chethiya.viharagama1@yaalalabs.com",
		"phone":     "+94768680618",
		"age":       29,
		"status":    "InActive",
	}

	reqTwo := helper.NewJSONRequest(t, http.MethodPost, "/users", reqBodyTwo)
	wTwo := httptest.NewRecorder()

	TestRouter.ServeHTTP(wTwo, reqTwo)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var createdUserTwo user.User
	errTwo := json.NewDecoder(wTwo.Body).Decode(&createdUserTwo)
	require.NoError(t, errTwo, "Failed to decode response")

	deleteReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", createdUserOne.UserId.String()), nil)
	deleteReq.Header.Set("Content-Type", "application/json")
	deleteW := httptest.NewRecorder()
	TestRouter.ServeHTTP(deleteW, deleteReq)

	//=========================================================

	reqAll := httptest.NewRequest(http.MethodGet, "/users", nil)
	wAll := httptest.NewRecorder()
	TestRouter.ServeHTTP(wAll, reqAll)

	assert.Equal(t, http.StatusOK, wAll.Code, "Expected status 200 OK")

	var users []user.User
	errAll := json.NewDecoder(wAll.Body).Decode(&users)

	require.NoError(t, errAll, "Failed to decode response")

	assert.Greater(t, len(users), 0, "Should have at least one user")
}
