package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlatten_Encode(t *testing.T) {
	t.Parallel()

	// Mock the flatten entity.
	f := Flatten{
		Date:        "2024-04-15",
		ProjectID:   "4974",
		NumTxs:      5,
		TotalVolume: 0.6136203411678249,
	}

	// Encode the flatten entity.
	encoded := f.Encode()

	// Check the encoded record.
	require.Equal(t, []string{
		"2024-04-15",
		"4974",
		"5",
		"0.6136203411678249",
	}, encoded)
}

func TestFlatten_Decode(t *testing.T) {
	t.Parallel()

	// Mock the CSV record.
	record := []string{
		"2024-04-15",
		"4974",
		"5",
		"0.6136203411678249",
	}

	var f Flatten

	// Decode the record into a flatten entity.
	err := f.Decode(record)
	require.NoError(t, err)

	// Check the decoded record.
	require.Equal(t, "2024-04-15", f.Date)
	require.Equal(t, "4974", f.ProjectID)
	require.Equal(t, 5, f.NumTxs)
	require.InEpsilon(t, 0.6136203411678249, f.TotalVolume, 0)
}
