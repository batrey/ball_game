package file

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// ReadInputFile reads the player's data file and returns a slice of slices.
func ReadInputFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), ",")
		for i := range row {
			row[i] = strings.TrimSpace(row[i])
		}
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("file is empty or in incorrect format")
	}
	return data, nil
}
