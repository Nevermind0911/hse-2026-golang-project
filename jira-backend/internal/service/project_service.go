package service

import (
	"context"
	"fmt"

	"hse-2026-golang-project/internal/models"
	"hse-2026-golang-project/jira-backend/internal/repository"

	pb "hse-2026-golang-project/internal/proto/connector"
)

type ProjectService struct {
	repo *repository.ProjectRepository
	grpcClient pb.ConnectorServiceClient
}

func NewProjectService(repo *repository.ProjectRepository, client pb.ConnectorServiceClient) *ProjectService {
	return &ProjectService{
		repo: repo,
		grpcClient: client,
	}
}

func (s *ProjectService) GetAll(ctx context.Context) ([]models.Project, error) {
	return s.repo.GetAll(ctx)
}



func (s *ProjectService) Delete(ctx context.Context, id int64) error {
	req := &pb.DeleteProjectRequest{ProjectId: id}
	
	_, err := s.grpcClient.DeleteProject(ctx, req)
	if err != nil {
		return fmt.Errorf("Error deleting a project via the connector: %w", err)
	}

	return nil
}

func (s *ProjectService) Update(ctx context.Context, key string) error {
	req := &pb.UpdateProjectRequest{ProjectKey: key}

	_, err := s.grpcClient.UpdateProject(ctx, req)
	if err != nil {
		return fmt.Errorf("update project via connector: %w", err)
	}

	return nil
}