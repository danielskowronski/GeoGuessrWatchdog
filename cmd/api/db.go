package main

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func unwrapFloat8Default(v pgtype.Float8) float64 {
	if !v.Valid {
		return 0
	}
	return v.Float64
}
func unwrapInt8Default(v pgtype.Int8) int64 {
	if !v.Valid {
		return 0
	}
	return v.Int64
}
func unwrapTimestamptzDefault(v pgtype.Timestamptz) time.Time {
	if !v.Valid {
		return time.Time{}
	}
	return v.Time
}
