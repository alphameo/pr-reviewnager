// Package di provides dependency injection
package di

import (
	"context"
	"errors"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/app"
	"github.com/alphameo/pr-reviewnager/internal/domain"
)

type RepositoryContainer interface {
	UserRepository() domain.UserRepository
	TeamRepository() domain.TeamRepository
	PullRequestRepository() domain.PullRequestRepository
	Close(ctx context.Context) error
}

type ServiceContainer struct {
	UserService        app.UserService
	TeamService        app.TeamService
	PullRequestService app.PullRequestService
}

func NewServiceContainer(repositoryContainer RepositoryContainer) (*ServiceContainer, error) {
	if repositoryContainer == nil {
		return nil, errors.New("storage cannot be nil")
	}
	prDomainServ, err := domain.NewDefaultPullRequestDomainService(
		repositoryContainer.UserRepository(),
		repositoryContainer.PullRequestRepository(),
		repositoryContainer.TeamRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain pull request service: %w", err)
	}

	teamServ, err := app.NewDefaultTeamService(
		repositoryContainer.TeamRepository(),
		repositoryContainer.UserRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create team service: %w", err)
	}

	userServ, err := app.NewDefaultUserService(repositoryContainer.UserRepository())
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %w", err)
	}

	prServ, err := app.NewDefaultPullRequestService(
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
