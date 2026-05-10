package app

import (
	"github.com/gorilla/mux"

	"hse-2026-golang-project/jira-backend/internal/handler"
)

func NewRouter(
	projectHandler *handler.ProjectHandler,
	issueHandler *handler.IssueHandler,
	graphHandler *handler.GraphHandler,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/projects", projectHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/v1/projects/{id:[0-9]+}", projectHandler.Delete).Methods("DELETE")
	r.HandleFunc("/api/v1/projects/{project}/update", projectHandler.Update).Methods("POST")

	r.HandleFunc("/api/v1/issues", issueHandler.GetByProject).Methods("GET")
	r.HandleFunc("/api/v1/projects/{project}/issues", issueHandler.GetByProject).Methods("GET")

	r.HandleFunc("/api/v1/graph/make/{task:[0-9]+}", graphHandler.Make).Methods("POST")
	r.HandleFunc("/api/v1/graph/get/{task:[0-9]+}", graphHandler.Get).Methods("GET")
	r.HandleFunc("/api/v1/projects/{project}/graph/{task:[0-9]+}", graphHandler.Get).Methods("GET")
	r.HandleFunc("/api/v1/projects/{project}/graph/{task:[0-9]+}", graphHandler.Make).Methods("POST")
	r.HandleFunc("/api/v1/projects/{project}/graph/{task:[0-9]+}/make", graphHandler.Make).Methods("POST")
	r.HandleFunc("/api/v1/isAnalyzed", graphHandler.IsAnalyzed).Methods("GET")
	r.HandleFunc("/api/v1/projects/{project}/analyzed", graphHandler.IsAnalyzed).Methods("GET")

	return r
}
