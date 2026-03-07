package connector

import (
	"context"

	"github.com/sirupsen/logrus"

	"hse-2026-golang-project/internal/config"
	"hse-2026-golang-project/internal/db"
	pb "hse-2026-golang-project/internal/proto/connector"
)

type GRPCServer struct {
	pb.UnimplementedConnectorServiceServer
	storage *db.Storage
	client  *JiraClient
	cfg     config.ProgramSettings
	log     *logrus.Logger
}

func NewGRPCServer(storage *db.Storage, client *JiraClient, cfg config.ProgramSettings, log *logrus.Logger) *GRPCServer {
	return &GRPCServer{
		storage: storage,
		client:  client,
		cfg:     cfg,
		log:     log,
	}
}

func (s *GRPCServer) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {}
func (s *GRPCServer) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
}
func (s *GRPCServer) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
}
func (s *GRPCServer) GetProjects(ctx context.Context, req *pb.GetProjectsRequest) (*pb.GetProjectsResponse, error) {