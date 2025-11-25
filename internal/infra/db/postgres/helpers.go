package postgres

import (
	"time"

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
