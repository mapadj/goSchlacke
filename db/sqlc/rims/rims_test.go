package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpsertRimsV1(t *testing.T) {
	arg := RandomRimV1()

	rim, err := testQueries.UpsertRimsV1(context.Background(), arg)

	require.NoError(t, err)  // Test for Error
	require.NotEmpty(t, rim) // Test for Result not Nil

	// Test for return values as expected
	require.Equal(t, arg.Code, rim.Code)
	require.Equal(t, arg.Width, rim.Width)
	require.Equal(t, arg.Height, rim.Height)
	require.Equal(t, arg.OnePiece, rim.OnePiece)
	require.Equal(t, arg.Diameter, rim.Diameter)
	require.Equal(t, arg.Material, rim.Material)

}
