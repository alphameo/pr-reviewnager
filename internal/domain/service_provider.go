package domain

import (
	"errors"
	"fmt"

)

type ServiceProvider interface {
	PullRequestDomainService() PullRequestDomainService
}

type DefaultServiceProvider struct {
	prServ PullRequestDomainService
}

func NewDefaultServiceProvider(storage RepositoryProvider) (*DefaultServiceProvider, error) {
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

	return &DefaultServiceProvider{prServ: prDomainServ}, nil
}

func (p *DefaultServiceProvider) PullRequestDomainService() PullRequestDomainService {
	return p.prServ
}
