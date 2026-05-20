package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"hse-2026-golang-project/jira-backend/internal/service"
)

type GraphHandler struct {
	service *service.GraphService
}

func NewGraphHandler(s *service.GraphService) *GraphHandler {
	return &GraphHandler{service: s}
}

func (h *GraphHandler) Make(w http.ResponseWriter, r *http.Request) {
	project := projectKeyFromRequest(r)
	if project == "" {
		writeError(w, http.StatusBadRequest, "project key is required")
		return
	}

	task, err := parseTask(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task")
		return
	}

	err = h.service.Make(r.Context(), project, task)
	if errors.Is(err, service.ErrProjectNotFound) {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	if errors.Is(err, service.ErrUnsupportedTask) {
		writeError(w, http.StatusBadRequest, "unsupported task")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare graph data")
		return
	}

	if err := writeJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *GraphHandler) Get(w http.ResponseWriter, r *http.Request) {
	project := projectKeyFromRequest(r)
	if project == "" {
		writeError(w, http.StatusBadRequest, "project key is required")
		return
	}

	task, err := parseTask(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task")
		return
	}

	data, err := h.service.Get(r.Context(), project, task)
	if errors.Is(err, service.ErrProjectNotFound) {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	if errors.Is(err, service.ErrUnsupportedTask) {
		writeError(w, http.StatusBadRequest, "unsupported task")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load graph data")
		return
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *GraphHandler) IsAnalyzed(w http.ResponseWriter, r *http.Request) {
	project := projectKeyFromRequest(r)
	if project == "" {
		writeError(w, http.StatusBadRequest, "project key is required")
		return
	}

	if err := writeJSON(w, http.StatusOK, map[string]bool{"analyzed": h.service.IsAnalyzed(project)}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func parseTask(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["task"])
}
