package dto

import (
	"time"

	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestDTO struct {
	ID          v.ID
	Title       string
	AuthorID    v.ID
	Status      string
	MergedAt    *time.Time
	ReviewerIDs []v.ID
}
