package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"hse-2026-golang-project/internal/db"
	"hse-2026-golang-project/jira-backend/internal/service"
)

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(s *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

func (h *ProjectHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load projects")
		return
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if errors.Is(err, db.ErrNotFound) {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	key := projectKeyFromRequest(r)
	if key == "" {
		writeError(w, http.StatusBadRequest, "[project key is required]")
		return
	}

	if err := h.service.Update(r.Context(), key); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update project")
	}

	if err := writeJSON(w, http.StatusOK, map[string]string {
		"status": "ok",
		"project": key,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}