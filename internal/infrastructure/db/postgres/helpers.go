package postgres

import (
	"time"

	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
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

func PullRequestToEntity(dbPR *db.PullRequest) (*e.PullRequest, error) {
	status, err := v.NewPRStatusFromString(dbPR.Status)
	if err != nil {
		return nil, err
	}
	var mergedAt *time.Time
	if dbPR.MergedAt.Valid {
		t := TimeFromTimestamptz(dbPR.MergedAt)
		mergedAt = &t
	} else {
		mergedAt = nil
	}

	return e.NewExistingPullRequest(
		v.ID(dbPR.ID),
		dbPR.Title,
		v.ID(dbPR.AuthorID),
		TimeFromTimestamptz(dbPR.CreatedAt),
		status,
		mergedAt,
		nil,
	)
}
