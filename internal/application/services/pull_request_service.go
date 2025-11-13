package services

import (
	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestService interface {
	CreatePullRequest(pullReques dto.PullRequestDTO) error
	MarkPullRequestAsMergedByID(pullRequestID v.ID) error
	ReassignReviewerByID(reviewerID v.ID, pullRequestID v.ID) error
	findPullRequestsForReviewerWithUserID(userID v.ID) ([]dto.PullRequestDTO, error)
}
