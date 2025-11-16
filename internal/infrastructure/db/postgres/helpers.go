package postgres

import (
	"time"

	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
	db "github.com/alphameo/pr-reviewnager/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func TimestampFromTime(t time.Time) pgtype.Timestamp {
	var ts pgtype.Timestamp
	ts.Scan(t)
	return ts
}

func TimeFromTimestamp(ts pgtype.Timestamp) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}

func PullRequestToEntity(dbPR *db.PullRequest) (*e.PullRequest, error) {
	status, err := v.NewPRStatusFromString(dbPR.Status)
	if err != nil {
		return nil, err
	}
	var mergedAt *time.Time
	if dbPR.MergedAt.Valid {
		t := TimeFromTimestamp(dbPR.MergedAt)
		mergedAt = &t
	} else {
		mergedAt = nil
	}

	pr := e.NewExistingPullRequest(
		v.ID(dbPR.ID),
		dbPR.Title,
		v.ID(dbPR.AuthorID),
		TimeFromTimestamp(dbPR.CreatedAt),
		status,
		mergedAt,
	)

	return pr, nil
}
