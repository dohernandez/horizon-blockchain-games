package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

// LoadSampleData loads sample data from a CSV file.
func LoadSampleData(limit, offset int) ([][]string, error) {
	// Open the CSV file.
	file, err := os.Open("../resources/sample_data.csv")
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	defer file.Close() //nolint:errcheck

	// Create a new CSV reader.
	reader := csv.NewReader(file)

	cursor := 0

	// When offset is not -1, keep the header.
	if offset != -1 {
		for {
			// Read each record (line) from the CSV.
			_, err := reader.Read()
			if errors.Is(err, io.EOF) {
				break // End of file
			}

			if err != nil {
				return nil, fmt.Errorf("moving cursor: %w", err)
			}

			if cursor >= offset {
				break
			}
		}
	}

	if offset == -1 {
		offset = 0
	}

	records := make([][]string, 0, limit)

	for {
		if cursor >= offset+limit {
			break
		}

		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break // End of file
		}

		if err != nil {
			return nil, fmt.Errorf("reading record: %w", err)
		}

		records = append(records, record)

		cursor++
	}

	return records, nil
}
