package postgres

import (
	"context"
	"errors"
	"time"

	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
	db "github.com/alphameo/pr-reviewnager/internal/infrastructure/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PullRequestRepository struct {
	queries *db.Queries
	dbConn  *pgx.Conn
}

func NewPullRequestRepository(queries *db.Queries, databaseConnection *pgx.Conn) (*PullRequestRepository, error) {
	if queries == nil {
		return nil, errors.New("queries cannot be nil")
	}
	if databaseConnection == nil {
		return nil, errors.New("database connection cannot be nil")
	}
	r := PullRequestRepository{
		queries: queries,
		dbConn:  databaseConnection,
	}
	return &r, nil
}

func (r *PullRequestRepository) Create(pullRequest *e.PullRequest) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	err = qtx.CreatePullRequest(ctx, db.CreatePullRequestParams{
		ID:       uuid.UUID(pullRequest.ID()),
		Title:    pullRequest.Title(),
		AuthorID: uuid.UUID(pullRequest.AuthorID()),
		Status:   pullRequest.Status().String(),
		MergedAt: TimestamptzFromTime(*pullRequest.MergedAt()),
	})
	if err != nil {
		return err
	}

	for _, reviewerID := range pullRequest.ReviewerIDs() {
		err = qtx.CreatePullRequestReviewer(ctx, db.CreatePullRequestReviewerParams{
			PullRequestID: uuid.UUID(pullRequest.ID()),
			ReviewerID:    uuid.UUID(reviewerID),
		})
		if err != nil {
			return err
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *PullRequestRepository) FindByID(id v.ID) (*e.PullRequest, error) {
	ctx := context.Background()

	rows, err := r.queries.GetPullRequestWithReviewersByID(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

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
				t := TimeFromTimestamptz(row.MergedAt)
				mergedAt = &t
			}

			pr = e.NewExistingPullRequest(
				v.ID(prID),
				row.Title,
				v.ID(row.AuthorID),
				TimeFromTimestamptz(row.CreatedAt),
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

	if len(prs) == 0 {
		return nil, nil
	}
	if len(prs) > 1 {
		return nil, errors.New("unexpected multiple PRs returned for a single ID")
	}

	return prs[0], nil
}

func (r *PullRequestRepository) FindAll() ([]*e.PullRequest, error) {
	ctx := context.Background()

	rows, err := r.queries.GetPullRequestsWithReviewers(ctx)
	if err != nil {
		return nil, err
	}

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
				t := TimeFromTimestamptz(row.MergedAt)
				mergedAt = &t
			}

			pr = e.NewExistingPullRequest(
				v.ID(prID),
				row.Title,
				v.ID(row.AuthorID),
				TimeFromTimestamptz(row.CreatedAt),
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

func (r *PullRequestRepository) Update(pullRequest *e.PullRequest) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	var mergedAt pgtype.Timestamptz
	if pullRequest.MergedAt() == nil {
		mergedAt = pgtype.Timestamptz{Valid: false}
	} else {
		mergedAt = TimestamptzFromTime(*pullRequest.MergedAt())
	}

	err = qtx.UpdatePullRequest(ctx, db.UpdatePullRequestParams{
		ID:       uuid.UUID(pullRequest.ID()),
		Title:    pullRequest.Title(),
		AuthorID: uuid.UUID(pullRequest.AuthorID()),
		Status:   pullRequest.Status().String(),
		MergedAt: mergedAt,
	})
	if err != nil {
		return err
	}

	err = qtx.DeletePullRequestReviewersByPRID(ctx, uuid.UUID(pullRequest.ID()))
	if err != nil {
		return err
	}

	for _, id := range pullRequest.ReviewerIDs() {
		err := qtx.CreatePullRequestReviewer(
			ctx,
			db.CreatePullRequestReviewerParams{
				PullRequestID: uuid.UUID(pullRequest.ID()),
				ReviewerID:    uuid.UUID(id),
			})
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *PullRequestRepository) DeleteByID(id v.ID) error {
	ctx := context.Background()

	err := r.queries.DeletePullRequest(ctx, uuid.UUID(id))
	if err != nil {
		return err
	}
	return nil
}

func (r *PullRequestRepository) FindPullRequestsByReviewer(userID v.ID) ([]*e.PullRequest, error) {
	ctx := context.Background()

	rows, err := r.queries.GetPullRequestsWithReviewersByReviewerID(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, err
	}

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
				t := TimeFromTimestamptz(row.MergedAt)
				mergedAt = &t
			}

			pr = e.NewExistingPullRequest(
				v.ID(prID),
				row.Title,
				v.ID(row.AuthorID),
				TimeFromTimestamptz(row.CreatedAt),
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
