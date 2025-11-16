package instrument

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	httputils "user-management/internal/common/httputils"
	"user-management/internal/middleware"
	ws "user-management/internal/ws"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Service  *Service
	Validate *validator.Validate
	ConnMgr  *ws.ConnectionManager
}

func NewHandler(service *Service, validate *validator.Validate, connMgr *ws.ConnectionManager) *Handler {
	return &Handler{
		Service:  service,
		Validate: validate,
		ConnMgr:  connMgr,
	}
}

// CreateInstrument godoc
// @Summary Create a new instrument
// @Description Create a new instrument
// @Tags instruments
// @Accept  json
// @Produce  json
// @Success 200 {object} Instrument
// @Failure      400  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
// @Router /instruments [post]
func (h *Handler) CreateInstrument(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req Instrument

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid request", "error", r)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid request", r)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		details := httputils.ConvertValidationErrors(err)
		slog.Warn("Instrument update failed", "error", "Validation failed")
		httputils.WriteDetailedError(w, http.StatusBadRequest, "Validation failed", details, r)
		return
	}

	instrument, err := h.Service.CreateInstrument(r.Context(), &req)

	if err != nil {
		slog.Error("Failed to create instrument", "error", r)
		httputils.WriteError(w, http.StatusInternalServerError, "Failed to create instrument", r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(instrument)
}

// GetInstrumentById godoc
// @Summary Get instrument by id
// @Description Get instrument details by id
// @Tags instruments
// @Accept  json
// @Produce  json
// @Success 200 {object} Instrument
// @Param id path string true "Instrument ID"
// @Failure      400  {object}  httputils.ErrorResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
// @Router /instruments/{id} [get]
func (h *Handler) GetInstrumentById(w http.ResponseWriter, r *http.Request) {

	instrumentId, uuiderr := httputils.ParseUUIDFromURL(r, "id")
	if uuiderr != nil {
		http.Error(w, "Invalid instrument ID format", http.StatusBadRequest)
		return
	}

	instruments, err := h.Service.GetInstrumentById(r.Context(), instrumentId.String())
	if err != nil {
		slog.Warn(fmt.Sprintf("Instrument not found with id: %s", instrumentId))
		httputils.WriteError(w, http.StatusNotFound, "Instrument not found", r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(instruments)
}

// GetInstruments godoc
// @Summary Get all instruments
// @Description Get all instruments
// @Tags instruments
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} Instrument
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
// @Router /instruments [get]
func (h *Handler) GetInstruments(w http.ResponseWriter, r *http.Request) {

	page := r.Context().Value(middleware.PageKey).(int)
	limit := r.Context().Value(middleware.LimitKey).(int)
	offset := (page - 1) * limit

	instruments, err := h.Service.ListInstrumentsPaged(r.Context(), limit, offset)
	if err != nil {
		slog.Warn("Failed to fetch instruments")
		httputils.WriteError(w, http.StatusNotFound, "Failed to fetch instruments", r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(instruments)
}

// UpdateInstrumentById godoc
// @Summary Update instrument by id
// @Description Update an instrument by id
// @Tags instruments
// @Accept  json
// @Produce  json
// @Param id path string true "Instrument ID"
// @Success 200 {object} Instrument
// @Failure      400  {object}  httputils.ErrorResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
// @Router /instruments/{id} [patch]
func (h *Handler) UpdateInstrumentById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	instrumentId, uuiderr := httputils.ParseUUIDFromURL(r, "id")
	if uuiderr != nil {
		httputils.WriteError(w, http.StatusNotFound, "Invalid instrument ID format", r)
		return
	}

	var req InstrumentUpdateRequest
	if err := httputils.DecodeAndValidateRequest(r, &req, h.Validate); err != nil {
		details := httputils.ConvertValidationErrors(err)
		slog.Warn("Instrument update failed", "error", err)
		httputils.WriteDetailedError(w, http.StatusBadRequest, "Validation failed", details, r)
		return
	}

	updatedInstrument, err := h.Service.UpdateInstrument(r.Context(), instrumentId.String(), &req, h.ConnMgr)

	if err != nil {
		slog.Warn("Instrument update failed", "error", err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error(), r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedInstrument)
}

// DeleteInstrument godoc
// @Summary Delete instrument by id
// @Description Delete an existing instrument by id
// @Tags instruments
// @Accept  json
// @Produce  json
// @Param id path string true "Instrument ID"
// @Success 204
// @Router /instruments/{id} [delete]
func (h *Handler) DeleteInstrumentById(w http.ResponseWriter, r *http.Request) {

	instrumentId, uuiderr := httputils.ParseUUIDFromURL(r, "id")
	if uuiderr != nil {
		http.Error(w, "Invalid instrument ID format", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteInstrumentById(r.Context(), instrumentId.String())
	if err != nil {
		httputils.WriteError(w, http.StatusNotFound, "Failed to fetch instruments", r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
