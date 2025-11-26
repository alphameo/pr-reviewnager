package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/alphameo/pr-reviewnager/internal/domain"
	db "github.com/alphameo/pr-reviewnager/internal/infra/db/sqlc"
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

	return &PullRequestRepository{
		queries: queries,
		dbConn:  databaseConnection,
	}, nil
}

func (r *PullRequestRepository) Create(pullRequest *domain.PullRequest) error {
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

	err = qtx.CreatePullRequest(ctx, db.CreatePullRequestParams{
		ID:       pullRequest.ID().Value(),
		Title:    pullRequest.Title().Value(),
		AuthorID: pullRequest.AuthorID().Value(),
		Status:   pullRequest.Status().String(),
		MergedAt: mergedAt,
	})
	if err != nil {
		return err
	}

	for _, reviewerID := range pullRequest.ReviewerIDs() {
		err = qtx.CreatePullRequestReviewer(ctx, db.CreatePullRequestReviewerParams{
			PullRequestID: pullRequest.ID().Value(),
			ReviewerID:    reviewerID.Value(),
		})
		if err != nil {
			return err
		}

	}

	return tx.Commit(ctx)
}

func (r *PullRequestRepository) FindByID(id domain.ID) (*domain.PullRequest, error) {
	ctx := context.Background()

	rows, err := r.queries.GetPullRequestWithReviewersByID(ctx, id.Value())
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	var reviewerIDs []domain.ID

	for _, row := range rows {
		if row.ReviewerID.Valid {
			reviewerID, err := domain.ParseID(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			reviewerIDs = append(reviewerIDs, reviewerID)
		}
	}

	var mergedAt *time.Time
	if rows[0].MergedAt.Valid {
		t := TimeFromTimestamptz(rows[0].MergedAt)
		mergedAt = &t
	}

	return domain.ExistingPullRequest(
		id,
		domain.ExistingPRTitle(rows[0].Title),
		domain.ID(rows[0].AuthorID),
		TimeFromTimestamptz(rows[0].CreatedAt),
		domain.ExistingPRStatus(rows[0].Status),
		mergedAt,
		reviewerIDs,
	), nil
}

func (r *PullRequestRepository) FindAll() ([]*domain.PullRequest, error) {
	ctx := context.Background()
	rows, err := r.queries.GetPullRequestsWithReviewers(ctx)
	if err != nil {
		return nil, err
	}

	type prData struct {
		Title     string
		AuthorID  uuid.UUID
		CreatedAt time.Time
		Status    string
		MergedAt  *time.Time
		Reviewers []domain.ID
	}
	prMap := make(map[uuid.UUID]*prData)

	for _, row := range rows {
		prID := row.ID

		if prMap[row.ID] == nil {
			var mergedAt *time.Time
			if row.MergedAt.Valid {
				t := TimeFromTimestamptz(row.MergedAt)
				mergedAt = &t
			}
			prMap[prID] = &prData{
				Title:     row.Title,
				AuthorID:  row.AuthorID,
				CreatedAt: TimeFromTimestamptz(row.CreatedAt),
				Status:    row.Status,
				MergedAt:  mergedAt,
				Reviewers: nil,
			}
		}

		if row.ReviewerID.Valid {
			reviewerID, err := domain.ParseID(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			prMap[prID].Reviewers = append(prMap[prID].Reviewers, reviewerID)
		}
	}

	prs := make([]*domain.PullRequest, 0, len(prMap))
	for id, data := range prMap {
		pr := domain.ExistingPullRequest(
			domain.ExistingID(id),
			domain.ExistingPRTitle(data.Title),
			domain.ExistingID(data.AuthorID),
			data.CreatedAt,
			domain.ExistingPRStatus(data.Status),
			data.MergedAt,
			data.Reviewers,
		)
		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *PullRequestRepository) Update(pullRequest *domain.PullRequest) error {
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
		ID:       pullRequest.ID().Value(),
		Title:    pullRequest.Title().Value(),
		AuthorID: pullRequest.AuthorID().Value(),
		Status:   pullRequest.Status().String(),
		MergedAt: mergedAt,
	})
	if err != nil {
		return err
	}

	err = qtx.DeletePullRequestReviewersByPRID(ctx, pullRequest.ID().Value())
	if err != nil {
		return err
	}

	for _, id := range pullRequest.ReviewerIDs() {
		err := qtx.CreatePullRequestReviewer(
			ctx,
			db.CreatePullRequestReviewerParams{
				PullRequestID: pullRequest.ID().Value(),
				ReviewerID:    id.Value(),
			})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PullRequestRepository) DeleteByID(id domain.ID) error {
	ctx := context.Background()

	err := r.queries.DeletePullRequest(ctx, id.Value())
	return err
}

func (r *PullRequestRepository) FindPullRequestsByReviewer(userID domain.ID) ([]*domain.PullRequest, error) {
	ctx := context.Background()
	rows, err := r.queries.GetPullRequestsWithReviewersByReviewerID(ctx, userID.Value())
	if err != nil {
		return nil, err
	}

	type prData struct {
		Title     string
		AuthorID  uuid.UUID
		CreatedAt time.Time
		Status    string
		MergedAt  *time.Time
		Reviewers []domain.ID
	}
	prMap := make(map[uuid.UUID]*prData)

	for _, row := range rows {
		prID := row.ID

		if prMap[row.ID] == nil {
			var mergedAt *time.Time
			if row.MergedAt.Valid {
				t := TimeFromTimestamptz(row.MergedAt)
				mergedAt = &t
			}
			prMap[prID] = &prData{
				Title:     row.Title,
				AuthorID:  row.AuthorID,
				CreatedAt: TimeFromTimestamptz(row.CreatedAt),
				Status:    row.Status,
				MergedAt:  mergedAt,
				Reviewers: nil,
			}
		}

		if row.ReviewerID.Valid {
			reviewerID, err := domain.ParseID(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			prMap[prID].Reviewers = append(prMap[prID].Reviewers, reviewerID)
		}
	}

	prs := make([]*domain.PullRequest, 0, len(prMap))
	for id, data := range prMap {
		pr := domain.ExistingPullRequest(
			domain.ExistingID(id),
			domain.ExistingPRTitle(data.Title),
			domain.ExistingID(data.AuthorID),
			data.CreatedAt,
			domain.ExistingPRStatus(data.Status),
			data.MergedAt,
			data.Reviewers,
		)
		prs = append(prs, pr)
	}

	return prs, nil
}
