package entities

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	// InputFieldNum is the number of fields in the input data. It is used to validate the input data before normalizing.
	InputFieldNum = 16
	// TSLayout is the layout of the timestamp use to parse the input data.
	TSLayout = "2006-01-02 15:04:05.000"
)

// Transaction represents a transaction entity.
type Transaction struct {
	TS                   time.Time
	Event                string
	ProjectID            string
	CurrencySymbol       string
	CurrencyValueDecimal float64
}

// TransactionNormalize normalizes the input data into a transaction entity.
func TransactionNormalize(d []string) (Transaction, error) {
	var t Transaction

	if len(d) != InputFieldNum {
		return t, fmt.Errorf("not enough fields in input: %d", len(d))
	}

	// Parse date.
	tsStr := d[1]

	ts, err := time.Parse(TSLayout, tsStr)
	if err != nil {
		return t, fmt.Errorf("parsing time: %w", err)
	}

	t.TS = ts
	t.Event = d[2]
	t.ProjectID = d[3]

	// Parse currency symbol.
	type propsJSON struct {
		CurrencySymbol string `json:"currencySymbol"`
	}

	var props propsJSON

	if err = json.Unmarshal([]byte(d[14]), &props); err != nil {
		return t, fmt.Errorf("parsing currency symbol: %w", err)
	}

	t.CurrencySymbol = props.CurrencySymbol

	// Parse currency value decimal.
	type valueJSON struct {
		CurrencyValueDecimal string `json:"currencyValueDecimal"`
	}

	var value valueJSON

	if err = json.Unmarshal([]byte(d[15]), &value); err != nil {
		return t, fmt.Errorf("parsing currency value decimal: %w", err)
	}

	// Convert string to float64.
	t.CurrencyValueDecimal, err = strconv.ParseFloat(value.CurrencyValueDecimal, 64)
	if err != nil {
		return t, fmt.Errorf("parsing currency value decimal")
	}

	return t, nil
}

// Encode encodes the transaction entity into a slice of strings.
func (t Transaction) Encode() []string {
	return []string{
		t.TS.Format(TSLayout),
		t.Event,
		t.ProjectID,
		t.CurrencySymbol,
		strconv.FormatFloat(t.CurrencyValueDecimal, 'g', -1, 64),
	}
}

// Decode decodes the transaction entity from a slice of strings.
func (t *Transaction) Decode(d []string) error {
	// Parse date.
	tsStr := d[0]

	ts, err := time.Parse(TSLayout, tsStr)
	if err != nil {
		return fmt.Errorf("parsing time: %w", err)
	}

	t.TS = ts
	t.Event = d[1]
	t.ProjectID = d[2]
	t.CurrencySymbol = d[3]

	// Convert string to float64.
	t.CurrencyValueDecimal, err = strconv.ParseFloat(d[4], 64)
	if err != nil {
		return fmt.Errorf("parsing currency value decimal")
	}

	return nil
}
