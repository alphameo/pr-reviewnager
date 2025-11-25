// Package di provides dependency injection
package di

import (
	"context"
	"errors"
	"fmt"

	s "github.com/alphameo/pr-reviewnager/internal/application/services"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	ds "github.com/alphameo/pr-reviewnager/internal/domain/services"
)

type RepositoryContainer interface {
	UserRepository() r.UserRepository
	TeamRepository() r.TeamRepository
	PullRequestRepository() r.PullRequestRepository
	Close(ctx context.Context) error
}

type ServiceContainer struct {
	UserService        s.UserService
	TeamService        s.TeamService
	PullRequestService s.PullRequestService
}

func NewServiceContainer(repositoryContainer RepositoryContainer) (*ServiceContainer, error) {
	if repositoryContainer == nil {
		return nil, errors.New("storage cannot be nil")
	}
	prDomainServ, err := ds.NewDefaultPullRequestDomainService(
		repositoryContainer.UserRepository(),
		repositoryContainer.PullRequestRepository(),
		repositoryContainer.TeamRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain pull request service: %w", err)
	}

	teamServ, err := s.NewDefaultTeamService(
		repositoryContainer.TeamRepository(),
		repositoryContainer.UserRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create team service: %w", err)
	}

	userServ, err := s.NewDefaultUserService(repositoryContainer.UserRepository())
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %w", err)
	}

	prServ, err := s.NewDefaultPullRequestService(
		prDomainServ,
		repositoryContainer.PullRequestRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request service: %w", err)
	}

	return &ServiceContainer{
		TeamService:        teamServ,
		UserService:        userServ,
		PullRequestService: prServ,
	}, nil
}
