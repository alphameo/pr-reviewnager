package postgres

import (
	"time"

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
