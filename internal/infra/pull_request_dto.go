package infra

import (
	"time"

	"github.com/alphameo/pr-reviewnager/internal/domain"
)

type PullRequestDTO struct {
	ID          domain.ID
	Title       string
	AuthorID    domain.ID
	CreatedAt   time.Time
	Status      string
	MergedAt    *time.Time
	ReviewerIDs []domain.ID
}

func (prDTO *PullRequestDTO) ToEntity() *domain.PullRequest {
	return domain.ExistingPullRequest(
		prDTO.ID,
		prDTO.Title,
		prDTO.AuthorID,
		prDTO.CreatedAt,
		domain.ExistingPRStatus(prDTO.Status),
		prDTO.MergedAt,
		prDTO.ReviewerIDs,
	)
}
