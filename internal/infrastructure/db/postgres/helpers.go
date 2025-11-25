package postgres

import (
	"time"

	"github.com/alphameo/pr-reviewnager/internal/domain"
	db "github.com/alphameo/pr-reviewnager/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func TimestamptzFromTime(t time.Time) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	ts.Scan(t)

	return ts
}

func TimeFromTimestamptz(ts pgtype.Timestamptz) time.Time {
	if ts.Valid {
		return ts.Time
	}

	return time.Time{}
}

func PullRequestToEntity(dbPR *db.PullRequest) (*domain.PullRequestDTO, error) {
	var mergedAt *time.Time
	if dbPR.MergedAt.Valid {
		t := TimeFromTimestamptz(dbPR.MergedAt)
		mergedAt = &t
	} else {
		mergedAt = nil
	}

	return &domain.PullRequestDTO{
		ID:          domain.ID(dbPR.ID),
		Title:       dbPR.Title,
		AuthorID:    domain.ID(dbPR.AuthorID),
		CreatedAt:   TimeFromTimestamptz(dbPR.CreatedAt),
		Status:      dbPR.Status,
		MergedAt:    mergedAt,
		ReviewerIDs: nil,
	}, nil
}
