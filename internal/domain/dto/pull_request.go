package dto

import (
	"time"

	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequest struct {
	ID          v.ID
	Title       string
	AuthorID    v.ID
	CreatedAt   time.Time
	Status      string
	MergedAt    *time.Time
	ReviewerIDs []v.ID
}
