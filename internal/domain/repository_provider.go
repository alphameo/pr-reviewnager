package domain

import "context"

type RepositoryProvider interface {
	UserRepository() UserRepository
	TeamRepository() TeamRepository
	PullRequestRepository() PullRequestRepository
	Close(ctx context.Context) error
}
