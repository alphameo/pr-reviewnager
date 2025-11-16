package postgres

import (
	"time"

	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
	db "github.com/alphameo/pr-reviewnager/internal/infrastructure/db/sqlc"
	"github.com/google/uuid"
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

	pr := e.NewPullRequestWithID(
		v.ID(dbPR.ID),
		dbPR.Title,
		v.ID(dbPR.AuthorID),
		status,
		mergedAt,
	)

	return pr, nil
}

func aggregatePRsFromRows(rows []db.GetPullRequestsWithReviewersRow) ([]*e.PullRequest, error) {
	prMap := make(map[uuid.UUID]*e.PullRequest)

	for _, row := range rows {
		prID := uuid.UUID(row.ID)
		pr, exists := prMap[prID]

		if !exists {
			status, err := v.NewPRStatusFromString(row.Status)
			if err != nil {
				return nil, err
			}
			var mergedAt *time.Time
			if row.MergedAt.Valid {
				t := TimeFromTimestamp(row.MergedAt)
				mergedAt = &t
			}

			pr = e.NewPullRequestWithID(
				v.ID(prID),
				row.Title,
				v.ID(row.AuthorID),
				status,
				mergedAt,
			)
			prMap[prID] = pr
		}

		if row.ReviewerID.Valid {
			reviewerID, err := v.NewIDFromString(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			err = pr.AssignReviewer(reviewerID)
			if err != nil {
				return nil, err
			}
		}
	}

	prs := make([]*e.PullRequest, 0, len(prMap))
	for _, pr := range prMap {
		prs = append(prs, pr)
	}

	return prs, nil
}

func aggregatePRsFromReviewerRows(rows []db.GetPullRequestsWithReviewersByReviewerIDRow) ([]*e.PullRequest, error) {
	prMap := make(map[uuid.UUID]*e.PullRequest)

	for _, row := range rows {
		prID := uuid.UUID(row.ID)
		pr, exists := prMap[prID]

		if !exists {
			status, err := v.NewPRStatusFromString(row.Status)
			if err != nil {
				return nil, err
			}
			var mergedAt *time.Time
			if row.MergedAt.Valid {
				t := TimeFromTimestamp(row.MergedAt)
				mergedAt = &t
			}

			pr = e.NewPullRequestWithID(
				v.ID(prID),
				row.Title,
				v.ID(row.AuthorID),
				status,
				mergedAt,
			)
			prMap[prID] = pr
		}

		if row.ReviewerID.Valid {
			reviewerID, err := v.NewIDFromString(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			err = pr.AssignReviewer(reviewerID)
			if err != nil {
				return nil, err
			}
		}
	}

	prs := make([]*e.PullRequest, 0, len(prMap))
	for _, pr := range prMap {
		prs = append(prs, pr)
	}

	return prs, nil
}

func aggregatePRsFromSinglePRRows(rows []db.GetPullRequestWithReviewersByIDRow) ([]*e.PullRequest, error) {
	prMap := make(map[uuid.UUID]*e.PullRequest)

	for _, row := range rows {
		prID := uuid.UUID(row.ID)
		pr, exists := prMap[prID]

		if !exists {
			status, err := v.NewPRStatusFromString(row.Status)
			if err != nil {
				return nil, err
			}
			var mergedAt *time.Time
			if row.MergedAt.Valid {
				t := TimeFromTimestamp(row.MergedAt)
				mergedAt = &t
			}

			pr = e.NewPullRequestWithID(
				v.ID(prID),
				row.Title,
				v.ID(row.AuthorID),
				status,
				mergedAt,
			)
			prMap[prID] = pr
		}

		if row.ReviewerID.Valid {
			reviewerID, err := v.NewIDFromString(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			err = pr.AssignReviewer(reviewerID)
			if err != nil {
				return nil, err
			}
		}
	}

	prs := make([]*e.PullRequest, 0, len(prMap))
	for _, pr := range prMap {
		prs = append(prs, pr)
	}

	return prs, nil
}
