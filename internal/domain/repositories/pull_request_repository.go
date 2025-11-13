package repositories

import (
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestRepository interface {
	Repository[e.PullRequest, v.ID]
	findPullRequestsByUserID(userID v.ID) ([]e.PullRequest, error)
}
