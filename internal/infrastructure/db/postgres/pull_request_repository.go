package postgres

import (
	"context"
	"errors"

	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
	db "github.com/alphameo/pr-reviewnager/internal/infrastructure/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PullRequestRepository struct {
	queries  *db.Queries
	database *pgx.Conn
}

func NewPullRequestRepository(queries *db.Queries) (*PullRequestRepository, error) {
	if queries != nil {
		return nil, errors.New("queries cannot be nil")
	}
	r := PullRequestRepository{queries: queries}
	return &r, nil
}

func (r *PullRequestRepository) Create(pullRequest *e.PullRequest) error {
	ctx := context.Background()
	tx, err := r.database.Begin(ctx)
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
		MergedAt: TimestampFromTime(*pullRequest.MergedAt()),
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
	tx, err := r.database.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	dbPR, err := qtx.GetPullRequest(ctx, uuid.UUID(id))
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	reviewerIDs, err := qtx.GetPullRequestReviewerReviewerIDs(ctx, uuid.UUID(id))
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	pr, err := PullRequestToEntity(&dbPR)
	if err != nil {
		return nil, err
	}
	for _, reviewerID := range reviewerIDs {
		err = pr.AssignReviewer(v.ID(reviewerID))
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (r *PullRequestRepository) FindAll() ([]*e.PullRequest, error) {
	ctx := context.Background()
	tx, err := r.database.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	dbPRs, err := qtx.GetPullRequests(ctx)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	prs := make([]*e.PullRequest, len(dbPRs))
	for i, dbPR := range dbPRs {
		reviewerIDs, err := qtx.GetPullRequestReviewerReviewerIDs(ctx, uuid.UUID(dbPR.ID))
		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		}

		pr, err := PullRequestToEntity(&dbPR)
		if err != nil {
			return nil, err
		}
		for _, reviewerID := range reviewerIDs {
			err = pr.AssignReviewer(v.ID(reviewerID))
			if err != nil {
				return nil, err
			}
		}

		prs[i] = pr
	}

	return prs, nil
}

func (r *PullRequestRepository) Update(pullRequest *e.PullRequest) error {
	ctx := context.Background()
	tx, err := r.database.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	var mergedAt pgtype.Timestamp
	if pullRequest.MergedAt() == nil {
		mergedAt = pgtype.Timestamp{Valid: false}
	} else {
		mergedAt = TimestampFromTime(*pullRequest.MergedAt())
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
		return nil
	}
	return nil
}

func (r *PullRequestRepository) FindPullRequestsByReviewer(userID v.ID) ([]*e.PullRequest, error) {
	ctx := context.Background()
	tx, err := r.database.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	dbPRs, err := qtx.GetPullRequestsByReviewer(ctx, uuid.UUID(userID))
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	prs := make([]*e.PullRequest, len(dbPRs))
	for i, dbPR := range dbPRs {
		reviewerIDs, err := qtx.GetPullRequestReviewerReviewerIDs(ctx, uuid.UUID(dbPR.ID))
		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		}

		pr, err := PullRequestToEntity(&dbPR)
		if err != nil {
			return nil, err
		}
		for _, reviewerID := range reviewerIDs {
			err = pr.AssignReviewer(v.ID(reviewerID))
			if err != nil {
				return nil, err
			}
		}

		prs[i] = pr
	}

	return prs, nil
}
