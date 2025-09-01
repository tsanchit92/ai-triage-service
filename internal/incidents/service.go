package incidents

import (
	"context"
	"errors"
	"strings"

	"github.com/yourname/ai-triage/internal/ai"
)

type Service struct {
	AI  ai.Client
	Repo Repository
}

func NewService(aiClient ai.Client, repo Repository) *Service {
	return &Service{AI: aiClient, Repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateIncidentRequest) (Incident, error) {
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Description) == "" || strings.TrimSpace(req.AffectedService) == "" {
		return Incident{}, errors.New("missing required fields")
	}
	class, err := s.AI.Classify(ctx, req.Title, req.Description, req.AffectedService)
	if err != nil {
		return Incident{}, err
	}
	inc := NewIncident(req.Title, req.Description, req.AffectedService, class.Severity, class.Category)
	if err := s.Repo.Insert(ctx, inc); err != nil {
		return Incident{}, err
	}
	return inc, nil
}

func (s *Service) List(ctx context.Context) ([]Incident, error) {
	return s.Repo.List(ctx)
}
