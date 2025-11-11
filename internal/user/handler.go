package user

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	httputils "user-management/internal/common/httputils"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(service *Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} User
// @Failure      400  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid request", "error", r)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid request", r)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		details := httputils.ConvertValidationErrors(err)
		slog.Warn("User update failed", "error", "Validation failed")
		httputils.WriteDetailedError(w, http.StatusBadRequest, "Validation failed", details, r)
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)

	if err != nil {
		slog.Error("Failed to create user", "error", r)
		httputils.WriteError(w, http.StatusInternalServerError, "Failed to create user", r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUserById godoc
// @Summary Get user by id
// @Description Get user details by id
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} User
// @Failure      400  {object}  httputils.ErrorResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
// @Router /users/{id} [get]
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {

	userId, uuiderr := httputils.ParseUUIDFromURL(r, "id")
	if uuiderr != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	users, err := h.service.GetUserById(r.Context(), userId.String())
	if err != nil {
		slog.Warn(fmt.Sprintf("User not found with id: %s", userId))
		httputils.WriteError(w, http.StatusNotFound, "User not found", r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
// @Router /users [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		slog.Warn("Failed to fetch users")
		httputils.WriteError(w, http.StatusNotFound, "Failed to fetch users", r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// UpdateUserById godoc
// @Summary Update user
// @Description Update an user by id
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} User
// @Failure      400  {object}  httputils.ErrorResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
// @Router /users/{id} [patch]
func (h *Handler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId, uuiderr := httputils.ParseUUIDFromURL(r, "id")
	if uuiderr != nil {
		http.Error(w, "Invalid user ID format", http.StatusNotFound)
		return
	}

	var req UserUpdateRequest
	if err := httputils.DecodeAndValidateRequest(r, &req, h.validate); err != nil {
		slog.Warn("User update failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), userId.String(), &req)

	if err != nil {
		slog.Warn("User update failed", "error", err)
		httputils.WriteError(w, http.StatusInternalServerError, "User update failed", r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete an existing user by id
// @Tags users
// @Accept  json
// @Produce  json
// @Success 204
// @Router /users/{id} [delete]
func (h *Handler) DeleteUserById(w http.ResponseWriter, r *http.Request) {

	userId, uuiderr := httputils.ParseUUIDFromURL(r, "id")
	if uuiderr != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteUserById(r.Context(), userId.String())
	if err != nil {
		httputils.WriteError(w, http.StatusNotFound, "Failed to fetch users", r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
