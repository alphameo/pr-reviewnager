// Package postgres provides repositories for PostgreSQL
package postgres

import (
	"context"

	db "github.com/alphameo/pr-reviewnager/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5"
)

func NewConnection(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewQueries(conn *pgx.Conn) *db.Queries {
	return db.New(conn)
}
