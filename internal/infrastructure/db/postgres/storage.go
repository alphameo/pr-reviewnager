package postgres

import (
	"context"
	"fmt"

	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	"github.com/jackc/pgx/v5"
)

type PSQLStorage struct {
	userRepo *UserRepository
	teamRepo *TeamRepository
	prRepo   *PullRequestRepository
	conn     *pgx.Conn
}

func NewPSQLStorage(ctx context.Context, dsn string) (*PSQLStorage, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	conn, err := NewConnection(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	queries := NewQueries(conn)

	teamRepo, err := NewTeamRepository(queries, conn)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create team repository: %w", err)
	}

	userRepo, err := NewUserRepository(queries)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	prRepo, err := NewPullRequestRepository(queries, conn)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create pull request repository: %w", err)
	}

	s := PSQLStorage{
		teamRepo: teamRepo,
		userRepo: userRepo,
		prRepo:   prRepo,
		conn:     conn,
	}

	return &s, nil
}

func (s *PSQLStorage) UserRepository() r.UserRepository {
	return s.userRepo
}

func (s *PSQLStorage) TeamRepository() r.TeamRepository {
	return s.teamRepo
}

func (s *PSQLStorage) PullRequestRepository() r.PullRequestRepository {
	return s.prRepo
}

func (s *PSQLStorage) Close(ctx context.Context) error {
	if s.conn == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return s.conn.Close(ctx)
}
