package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"hse-2026-golang-project/internal/config"
	"hse-2026-golang-project/internal/db"
	connector "hse-2026-golang-project/internal/jira"
	pb "hse-2026-golang-project/internal/proto/connector"
)

func initDB(cfg config.DBSettings) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := database.Ping(); err != nil {
		return nil, err
	}
	return database, nil
}

func main() {
	log := connector.NewLogger()
	log.Info("Starting Jira Connector...")

	cfgPath := "configs/config.yaml"
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	writeDB, err := initDB(cfg.WriteDB)
	if err != nil {
		log.Fatalf("Failed to connect to Write DB: %v", err)
	}
	defer writeDB.Close()

	readDB, err := initDB(cfg.ReadDB)
	if err != nil {
		log.Fatalf("Failed to connect to Read DB: %v", err)
	}
	defer readDB.Close()

	storage := db.NewStorage(writeDB, readDB)

	jiraClient := connector.NewJiraClient(cfg.Program, log)
	grpcServerInstance := connector.NewGRPCServer(storage, jiraClient, cfg.Program, log)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			start := time.Now()

			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic in gRPC handler: %v", r)
					log.WithFields(map[string]interface{}{
						"method": info.FullMethod,
						"panic":  r,
					}).Error("CRITICAL PANIC RECOVERED")
				}
			}()

			resp, err = handler(ctx, req)

			log.WithFields(map[string]interface{}{
				"method":   info.FullMethod,
				"duration": time.Since(start).String(),
				"error":    err,
			}).Info("gRPC call")
			return resp, err
		}),
	)

	pb.RegisterConnectorServiceServer(server, grpcServerInstance)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthSrv)
	healthSrv.SetServingStatus("connector", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(server)

	port := fmt.Sprintf(":%d", cfg.Program.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Infof("gRPC server listening on port %s", port)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gRPC server...")
	server.GracefulStop()
	log.Info("Server stopped.")
}
