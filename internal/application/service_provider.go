// Package app provides application layer
package app

import (
	"errors"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/domain"
)

type ServiceProvider interface {
	UserService() UserService
	TeamService() TeamService
	PullRequestService() PullRequestService
}

type DefaultServiceProvider struct {
	userServ *DefaultUserService
	teamServ *DefaultTeamService
	prServ   *DefaultPullRequestService
}

func NewDefaultServiceProvider(storage domain.RepositoryProvider, domainServiceProvider domain.ServiceProvider) (*DefaultServiceProvider, error) {
	if storage == nil {
		return nil, errors.New("storage cannot be nil")
	}
	if domainServiceProvider == nil {
		return nil, errors.New("domainServiceProvider cannot be nil")
	}

	teamServ, err := NewDefaultTeamService(
		storage.TeamRepository(),
		storage.UserRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create team service: %w", err)
	}

	userServ, err := NewDefaultUserService(storage.UserRepository())
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %w", err)
	}

	prServ, err := NewDefaultPullRequestService(
		domainServiceProvider.PullRequestDomainService(),
		storage.PullRequestRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request service: %w", err)
	}

	return &DefaultServiceProvider{
		teamServ: teamServ,
		userServ: userServ,
		prServ:   prServ,
	}, nil
}

func (p *DefaultServiceProvider) UserService() UserService {
	return p.userServ
}

func (p *DefaultServiceProvider) TeamService() TeamService {
	return p.teamServ
}

func (p *DefaultServiceProvider) PullRequestService() PullRequestService {
	return p.prServ
}
