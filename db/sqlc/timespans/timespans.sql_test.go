package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpsertTimepansV1(t *testing.T) {

	date := time.Date(2019, 2, 13, 0, 0, 0, 0, time.FixedZone("", 0))
	date2 := date.Add(time.Hour * 24)

	arg := UpsertTimespansV1Params{
		SchwackeID:   1,
		SchwackeCode: sql.NullInt32{Int32: 123, Valid: true},
		ValidFrom:    sql.NullTime{Time: date, Valid: true},
		ValidUntil:   sql.NullTime{Time: date2, Valid: true},
	}

	timespan, err := testQueries.UpsertTimespansV1(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, timespan)

	require.Equal(t, arg.SchwackeID, timespan.SchwackeID)
	require.Equal(t, arg.SchwackeCode, timespan.SchwackeCode)
	require.Equal(t, arg.ValidFrom, timespan.ValidFrom)
	require.Equal(t, arg.ValidUntil, timespan.ValidUntil)

}
