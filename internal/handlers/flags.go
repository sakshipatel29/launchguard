package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sakshipatel29/launchguard/internal/models"
	"github.com/sakshipatel29/launchguard/internal/store"
)

type FeatureFlagHandler struct {
	store *store.FeatureFlagStore
}

func NewFeatureFlagHandler(flagStore *store.FeatureFlagStore) *FeatureFlagHandler {
	return &FeatureFlagHandler{
		store: flagStore,
	}
}

func (h *FeatureFlagHandler) CreateFlag(w http.ResponseWriter, r *http.Request) {
	var req models.CreateFeatureFlagRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Key == "" || req.Environment == "" {
		writeError(w, http.StatusBadRequest, "name, key, and environment are required")
		return
	}

	if req.RolloutPercentage < 0 || req.RolloutPercentage > 100 {
		writeError(w, http.StatusBadRequest, "rollout_percentage must be between 0 and 100")
		return
	}

	flag, err := h.store.Create(req)
	if err != nil {
		if errors.Is(err, store.ErrDuplicateFlagKey) {
			writeError(w, http.StatusConflict, "feature flag key already exists for this environment")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to create feature flag")
		return
	}

	writeJSON(w, http.StatusCreated, flag)
}

func (h *FeatureFlagHandler) ListFlags(w http.ResponseWriter, r *http.Request) {
	flags := h.store.List()
	writeJSON(w, http.StatusOK, flags)
}

func (h *FeatureFlagHandler) GetFlag(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	flag, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrFlagNotFound) {
			writeError(w, http.StatusNotFound, "feature flag not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to get feature flag")
		return
	}

	writeJSON(w, http.StatusOK, flag)
}

func (h *FeatureFlagHandler) UpdateFlag(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateFeatureFlagRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Environment == "" {
		writeError(w, http.StatusBadRequest, "name and environment are required")
		return
	}

	if req.RolloutPercentage < 0 || req.RolloutPercentage > 100 {
		writeError(w, http.StatusBadRequest, "rollout_percentage must be between 0 and 100")
		return
	}

	flag, err := h.store.Update(id, req)
	if err != nil {
		if errors.Is(err, store.ErrFlagNotFound) {
			writeError(w, http.StatusNotFound, "feature flag not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to update feature flag")
		return
	}

	writeJSON(w, http.StatusOK, flag)
}

func (h *FeatureFlagHandler) DeleteFlag(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.store.Delete(id)
	if err != nil {
		if errors.Is(err, store.ErrFlagNotFound) {
			writeError(w, http.StatusNotFound, "feature flag not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to delete feature flag")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "feature flag deleted successfully",
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{
		"error": message,
	})
}
