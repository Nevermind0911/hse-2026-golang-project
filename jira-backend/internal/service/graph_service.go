package service

import (
	"context"
	"errors"
	"sync"

	"hse-2026-golang-project/jira-backend/internal/repository"
)

var ErrUnsupportedTask = errors.New("unsupported graph task")

type GraphService struct {
	repo     *repository.ProjectRepository
	mu       sync.RWMutex
	analyzed map[string]bool
}

func NewGraphService(repo *repository.ProjectRepository) *GraphService {
	return &GraphService{
		repo:     repo,
		analyzed: make(map[string]bool),
	}
}

func (s *GraphService) Make(ctx context.Context, projectKey string, task int) error {
	if task != 1 {
		return ErrUnsupportedTask
	}

	project, err := s.repo.GetByKey(ctx, projectKey)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrProjectNotFound
	}

	s.mu.Lock()
	s.analyzed[projectKey] = true
	s.mu.Unlock()

	return nil
}

func (s *GraphService) Get(ctx context.Context, projectKey string, task int) (map[string]int, error) {
	if task != 1 {
		return nil, ErrUnsupportedTask
	}

	project, err := s.repo.GetByKey(ctx, projectKey)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	issues, err := s.repo.GetIssuesByProject(ctx, project.JiraID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int, len(issues))
	for _, issue := range issues {
		result[issue.Priority]++
	}

	return result, nil
}

func (s *GraphService) IsAnalyzed(projectKey string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.analyzed[projectKey]
}
