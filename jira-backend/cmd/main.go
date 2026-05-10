package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "hse-2026-golang-project/internal/proto/connector"

	"hse-2026-golang-project/internal/db"
	"hse-2026-golang-project/jira-backend/internal/app"
	"hse-2026-golang-project/jira-backend/internal/handler"
	"hse-2026-golang-project/jira-backend/internal/repository"
	"hse-2026-golang-project/jira-backend/internal/service"
)

func main() {
	dsn := "postgres://pguser:pgpwd@localhost:5432/testdb?sslmode=disable"

	writeDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open master db: %v", err)
	}
	readDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open replica db: %v", err)
	}

	storage := db.NewStorage(writeDB, readDB)
	defer func() {
		if err := storage.Close(); err != nil {
			log.Printf("close db connections: %v", err)
		}
	}()

	repo := repository.NewProjectRepository(storage)

	connectorAddress := "connector:8001" 
	conn, err := grpc.NewClient(connectorAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("create grpc connection: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewConnectorServiceClient(conn)

	projectService := service.NewProjectService(repo, grpcClient)
	issueService := service.NewIssueService(repo)
	graphService := service.NewGraphService(repo)

	projectHandler := handler.NewProjectHandler(projectService)
	issueHandler := handler.NewIssueHandler(issueService)
	graphHandler := handler.NewGraphHandler(graphService)

	router := app.NewRouter(projectHandler, issueHandler, graphHandler)

	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
