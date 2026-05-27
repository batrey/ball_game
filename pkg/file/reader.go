package file

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

// Reader loads player visibility data from files.
type Reader struct{}

// ReadPlayers reads player visibility rows from the given file.
func (Reader) ReadPlayers(fileName string) ([][]string, error) {
	return ReadInputFile(fileName)
}

// ReadInputFile reads CSV player visibility rows from a file and ignores blank lines.
func ReadInputFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	var data [][]string
	for {
		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		for i := range row {
			row[i] = strings.TrimSpace(row[i])
		}
		if len(row) == 1 && row[0] == "" {
			continue
		}
		data = append(data, row)
	}

	if len(data) == 0 {
		return nil, errors.New("file is empty or in incorrect format")
	}
	return data, nil
}
