package repositories

import (
	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestRepository interface {
	Repository[e.PullRequest, dto.PullRequestDTO, v.ID]
	FindPullRequestsByReviewer(userID v.ID) ([]*dto.PullRequestDTO, error)
}
