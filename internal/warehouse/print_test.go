package warehouse

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

func TestPrint_Save(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Mock the flatten entity.
	flatten := entities.Flatten{
		Date:        "2024-04-15",
		ProjectID:   "4974",
		NumTxs:      5,
		TotalVolume: 0.6136203411678249,
	}

	// Redirect os.Stdout to capture output
	var buf bytes.Buffer

	originalStdout := os.Stdout
	r, w, _ := os.Pipe() //nolint:errcheck
	os.Stdout = w

	// Save the flatten entity.
	p := &Print{}

	err := p.Save(ctx, flatten)
	require.NoError(t, err)

	// Close the writer and restore os.Stdout
	w.Close() //nolint:errcheck,gosec

	os.Stdout = originalStdout

	// Copy the captured output to our buffer
	buf.ReadFrom(r) //nolint:errcheck,gosec

	// Check the printed record.
	require.Equal(t, "Save: {2024-04-15 4974 5 0.6136203411678249}\n", buf.String())
}
