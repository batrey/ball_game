package file

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// Reader loads player visibility data from files.
type Reader struct{}

// ReadPlayers reads player visibility rows from the given file.
func (Reader) ReadPlayers(fileName string) ([][]string, error) {
	return ReadInputFile(fileName)
}

// ReadInputFile reads player visibility rows from a file and ignores blank lines.
func ReadInputFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		row := strings.Split(line, ",")
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
