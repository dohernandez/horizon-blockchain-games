package entities

import (
	"fmt"
	"strconv"
)

// Flatten represents a flattened transaction entity.
type Flatten struct {
	Date        string  `bigquery:"date"`
	ProjectID   string  `bigquery:"project_id"`
	NumTxs      int     `bigquery:"num_transactions"`
	TotalVolume float64 `bigquery:"total_volume_usd"`
}

// Encode encodes the flatten entity into a slice of strings.
func (f Flatten) Encode() []string {
	return []string{
		f.Date,
		f.ProjectID,
		fmt.Sprintf("%d", f.NumTxs),
		strconv.FormatFloat(f.TotalVolume, 'g', -1, 64),
	}
}

// Decode decodes the flatten entity from a slice of strings.
func (f *Flatten) Decode(d []string) error {
	f.Date = d[0]
	f.ProjectID = d[1]

	var err error

	// Convert string to int.
	f.NumTxs, err = strconv.Atoi(d[2])
	if err != nil {
		return fmt.Errorf("parsing num txs")
	}

	// Convert string to float64.
	f.TotalVolume, err = strconv.ParseFloat(d[3], 64)
	if err != nil {
		return fmt.Errorf("parsing total volume")
	}

	return nil
}
