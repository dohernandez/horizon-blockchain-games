package entities_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

func TestTransactionNormalize(t *testing.T) {
	t.Parallel()

	// Mock the CSV record.
	record := []string{
		"seq-market",
		"2024-04-15 02:15:07.167",
		"BUY_ITEMS",
		"4974",
		"",
		"1",
		"0896ae95dcaeee38e83fa5c43bef99780d7b2be23bcab36214",
		"5d8afd8fec2fbf3e",
		"DE",
		"desktop",
		"linux",
		"x86_64",
		"chrome",
		"122.0.0.0",
		`{"tokenId":"215","txnHash":"0xd919290e80df271e77d1cbca61f350d2727531e0334266671ec20d626b2104a2","chainId":"137","collectionAddress":"0x22d5f9b75c524fec1d6619787e582644cd4d7422","currencyAddress":"0xd1f9c58e33933a993a3891f8acfe05a68e1afc05","currencySymbol":"SFL","marketplaceType":"amm","requestId":""}`,
		`{"currencyValueDecimal":"0.6136203411678249","currencyValueRaw":"613620341167824900"}`,
	}

	// Decode the record into a transaction.
	tx, err := entities.TransactionNormalize(record)
	require.NoError(t, err)

	ts, err := time.Parse("2006-01-02 15:04:05.000", "2024-04-15 02:15:07.167")
	require.NoError(t, err)

	require.Equal(t, ts, tx.TS)
	require.Equal(t, "BUY_ITEMS", tx.Event)
	require.Equal(t, "4974", tx.ProjectID)
	require.Equal(t, "SFL", tx.CurrencySymbol)
	require.InEpsilon(t, 0.6136203411678249, tx.CurrencyValueDecimal, 0)
}

func TestTransaction_Encode(t *testing.T) {
	t.Parallel()

	// Mock the transaction.
	tx := entities.Transaction{
		TS:                   time.Date(2024, 4, 15, 2, 15, 7, 167000000, time.UTC),
		Event:                "BUY_ITEMS",
		ProjectID:            "4974",
		CurrencySymbol:       "SFL",
		CurrencyValueDecimal: 0.6136203411678249,
	}

	// Encode the transaction.
	encoded := tx.Encode()

	// Check the encoded record.
	require.Equal(t, []string{
		time.Date(2024, 4, 15, 2, 15, 7, 167000000, time.UTC).Format("2006-01-02 15:04:05.000"),
		"BUY_ITEMS",
		"4974",
		"SFL",
		"0.6136203411678249",
	}, encoded)
}

func TestTransaction_Decode(t *testing.T) {
	t.Parallel()

	// Mock the CSV record.
	record := []string{
		"2024-04-15 02:15:07.167",
		"BUY_ITEMS",
		"4974",
		"SFL",
		"0.6136203411678249",
	}

	var tx entities.Transaction

	// Decode the record into a transaction.
	err := tx.Decode(record)
	require.NoError(t, err)

	// Check the decoded record.
	ts, err := time.Parse("2006-01-02 15:04:05.000", "2024-04-15 02:15:07.167")
	require.NoError(t, err)

	require.Equal(t, ts, tx.TS)
	require.Equal(t, "BUY_ITEMS", tx.Event)
	require.Equal(t, "4974", tx.ProjectID)
	require.Equal(t, "SFL", tx.CurrencySymbol)
	require.InEpsilon(t, 0.6136203411678249, tx.CurrencyValueDecimal, 0)
}
