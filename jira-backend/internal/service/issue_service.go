package service

import (
	"context"
	"errors"

	"hse-2026-golang-project/internal/models"
	"hse-2026-golang-project/jira-backend/internal/repository"
)

var ErrProjectNotFound = errors.New("project not found")

type IssueService struct {
	repo *repository.ProjectRepository
}

func NewIssueService(repo *repository.ProjectRepository) *IssueService {
	return &IssueService{repo: repo}
}

func (s *IssueService) GetByProjectKey(ctx context.Context, key string) ([]models.Issue, error) {
	project, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	return s.repo.GetIssuesByProject(ctx, project.JiraID)
}
