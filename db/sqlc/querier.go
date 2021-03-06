// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	CountRimsV1(ctx context.Context) (int64, error)
	CountTimespansV1(ctx context.Context) (int64, error)
	InsertLog(ctx context.Context, arg InsertLogParams) (Log, error)
	UpsertRimsV1(ctx context.Context, arg UpsertRimsV1Params) (Rim, error)
	UpsertTimespansV1(ctx context.Context, arg UpsertTimespansV1Params) (Timespan, error)
}

var _ Querier = (*Queries)(nil)
