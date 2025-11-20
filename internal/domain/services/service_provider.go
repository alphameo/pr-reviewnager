package services

import (
	"errors"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/infrastructure"
)

type ServiceProvider interface {
	PullRequestDomainService() PullRequestDomainService
}

type DefaultServiceProvider struct {
	prServ PullRequestDomainService
}

func NewDefaultServiceProvider(storage infrastructure.Storage) (*DefaultServiceProvider, error) {
	if storage == nil {
		return nil, errors.New("storage cannot be nil")
	}

	prDomainServ, err := NewDefaultPullRequestDomainService(
		storage.UserRepository(),
		storage.PullRequestRepository(),
		storage.TeamRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain pull request service: %w", err)
	}

	p := DefaultServiceProvider{prServ: prDomainServ}
	return &p, nil
}

func (p *DefaultServiceProvider) PullRequestDomainService() PullRequestDomainService {
	return p.prServ
}
