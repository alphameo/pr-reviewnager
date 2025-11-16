package dto

import (
	"time"

	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestDTO struct {
	ID          v.ID
	Title       string
	AuthorID    v.ID
	CreatedAt time.Time
	Status      string
	MergedAt    *time.Time
	ReviewerIDs []v.ID
}

type PullRequestWithNewReviewerIDDTO struct {
	PullRequest       *PullRequestDTO
	NewReviewerUserID v.ID
}
