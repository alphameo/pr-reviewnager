package app

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

type NewPullRequestDTO struct {
	ID       domain.ID
	Title    string
	AuthorID domain.ID
}
