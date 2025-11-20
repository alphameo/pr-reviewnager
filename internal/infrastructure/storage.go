// Package infrastructure provides objects, that interacts with infrastructure of application
package infrastructure

import (
	"context"

	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
)

type Storage interface {
	UserRepository() r.UserRepository
	TeamRepository() r.TeamRepository
	PullRequestRepository() r.PullRequestRepository
	Close(ctx context.Context) error
}
