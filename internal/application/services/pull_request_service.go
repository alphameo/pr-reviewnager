package services

import (
	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestService interface {
	CreatePullRequest(pullRequest *dto.PullRequestDTO) error
	MarkAsMerged(pullRequestID v.ID) (*dto.PullRequestDTO, error)
	ReassignReviewer(userID v.ID, pullRequestID v.ID) (*dto.PullRequestWithNewReviewerIDDTO, error)
	FindPullRequestsByReviewer(userID v.ID) ([]dto.PullRequestDTO, error)
}
