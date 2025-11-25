package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/alphameo/pr-reviewnager/internal/domain"
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
		ID:       uuid.UUID(pullRequest.ID()),
		Title:    pullRequest.Title(),
		AuthorID: uuid.UUID(pullRequest.AuthorID()),
		Status:   pullRequest.Status().String(),
		MergedAt: mergedAt,
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

	return tx.Commit(ctx)
}

func (r *PullRequestRepository) FindByID(id domain.ID) (*domain.PullRequestDTO, error) {
	ctx := context.Background()

	rows, err := r.queries.GetPullRequestWithReviewersByID(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	prMap := make(map[uuid.UUID]*domain.PullRequestDTO)

	for _, row := range rows {
		prID := uuid.UUID(row.ID)
		pr, exists := prMap[prID]

		if !exists {
			var mergedAt *time.Time
			if row.MergedAt.Valid {
				t := TimeFromTimestamptz(row.MergedAt)
				mergedAt = &t
			}

			pr = &domain.PullRequestDTO{
				ID:          domain.ID(prID),
				Title:       row.Title,
				AuthorID:    domain.ID(row.AuthorID),
				CreatedAt:   TimeFromTimestamptz(row.CreatedAt),
				Status:      row.Status,
				MergedAt:    mergedAt,
				ReviewerIDs: make([]domain.ID, 0),
			}
			prMap[prID] = pr
		}

		if row.ReviewerID.Valid {
			reviewerID, err := domain.NewIDFromString(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			pr.ReviewerIDs = append(prMap[prID].ReviewerIDs, reviewerID)
		}
	}

	prs := make([]*domain.PullRequestDTO, 0, len(prMap))
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

func (r *PullRequestRepository) FindAll() ([]*domain.PullRequestDTO, error) {
	ctx := context.Background()

	rows, err := r.queries.GetPullRequestsWithReviewers(ctx)
	if err != nil {
		return nil, err
	}

	prMap := make(map[uuid.UUID]*domain.PullRequestDTO)

	for _, row := range rows {
		prID := uuid.UUID(row.ID)
		pr, exists := prMap[prID]

		if !exists {
			var mergedAt *time.Time
			if row.MergedAt.Valid {
				t := TimeFromTimestamptz(row.MergedAt)
				mergedAt = &t
			}

			pr = &domain.PullRequestDTO{
				ID:          domain.ID(prID),
				Title:       row.Title,
				AuthorID:    domain.ID(row.AuthorID),
				CreatedAt:   TimeFromTimestamptz(row.CreatedAt),
				Status:      row.Status,
				MergedAt:    mergedAt,
				ReviewerIDs: make([]domain.ID, 0),
			}
			prMap[prID] = pr
		}

		if row.ReviewerID.Valid {
			reviewerID, err := domain.NewIDFromString(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			pr.ReviewerIDs = append(prMap[prID].ReviewerIDs, reviewerID)
		}
	}

	prs := make([]*domain.PullRequestDTO, 0, len(prMap))
	for _, pr := range prMap {
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

	return tx.Commit(ctx)
}

func (r *PullRequestRepository) DeleteByID(id domain.ID) error {
	ctx := context.Background()

	err := r.queries.DeletePullRequest(ctx, uuid.UUID(id))
	if err != nil {
		return err
	}

	return nil
}

func (r *PullRequestRepository) FindPullRequestsByReviewer(userID domain.ID) ([]*domain.PullRequestDTO, error) {
	ctx := context.Background()

	rows, err := r.queries.GetPullRequestsWithReviewersByReviewerID(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, err
	}

	prMap := make(map[uuid.UUID]*domain.PullRequestDTO)

	for _, row := range rows {
		prID := uuid.UUID(row.ID)
		pr, exists := prMap[prID]

		if !exists {
			var mergedAt *time.Time
			if row.MergedAt.Valid {
				t := TimeFromTimestamptz(row.MergedAt)
				mergedAt = &t
			}

			prMap[prID] = &domain.PullRequestDTO{
				ID:          domain.ID(prID),
				Title:       row.Title,
				AuthorID:    domain.ID(row.AuthorID),
				CreatedAt:   TimeFromTimestamptz(row.CreatedAt),
				Status:      row.Status,
				MergedAt:    mergedAt,
				ReviewerIDs: make([]domain.ID, 0),
			}
			prMap[prID] = pr
		}

		if row.ReviewerID.Valid {
			reviewerID, err := domain.NewIDFromString(row.ReviewerID.String())
			if err != nil {
				return nil, err
			}
			pr.ReviewerIDs = append(prMap[prID].ReviewerIDs, reviewerID)
		}
	}

	prs := make([]*domain.PullRequestDTO, 0, len(prMap))
	for _, pr := range prMap {
		prs = append(prs, pr)
	}

	return prs, nil
}
