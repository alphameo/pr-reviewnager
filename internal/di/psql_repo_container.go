package di

import (
	"context"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/domain"
	"github.com/alphameo/pr-reviewnager/internal/infra/db/postgres"
	"github.com/jackc/pgx/v5"
)

type PSQLRepositoryContainer struct {
	userRepo *postgres.UserRepository
	teamRepo *postgres.TeamRepository
	prRepo   *postgres.PullRequestRepository
	conn     *pgx.Conn
}

func NewPSQLRepositoryContainer(ctx context.Context, dsn string) (*PSQLRepositoryContainer, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	conn, err := postgres.NewConnection(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	queries := postgres.NewQueries(conn)

	teamRepo, err := postgres.NewTeamRepository(queries, conn)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create team repository: %w", err)
	}

	userRepo, err := postgres.NewUserRepository(queries)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	prRepo, err := postgres.NewPullRequestRepository(queries, conn)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create pull request repository: %w", err)
	}

	return &PSQLRepositoryContainer{
		teamRepo: teamRepo,
		userRepo: userRepo,
		prRepo:   prRepo,
		conn:     conn,
	}, nil
}

func (s *PSQLRepositoryContainer) UserRepository() domain.UserRepository {
	return s.userRepo
}

func (s *PSQLRepositoryContainer) TeamRepository() domain.TeamRepository {
	return s.teamRepo
}

func (s *PSQLRepositoryContainer) PullRequestRepository() domain.PullRequestRepository {
	return s.prRepo
}

func (s *PSQLRepositoryContainer) Close(ctx context.Context) error {
	if s.conn == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return s.conn.Close(ctx)
}
