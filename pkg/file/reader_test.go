package file

import (
	"os"
	"testing"
)

func TestReadInputFile(t *testing.T) {
	// Test cases for ReadInputFile function
	tests := []struct {
		name          string
		fileContent   string
		expectedData  [][]string
		expectedError string
	}{
		{
			name: "basic",
			fileContent: `George,Beth,Sue
		Rick,Anne
		Anne,Beth
		Beth,Anne,George
		Sue,Beth`,
			expectedData: [][]string{
				{"George", "Beth", "Sue"},
				{"Rick", "Anne"},
				{"Anne", "Beth"},
				{"Beth", "Anne", "George"},
				{"Sue", "Beth"},
			},
			expectedError: "",
		},
		{
			name:          "empty file",
			fileContent:   "",
			expectedData:  nil,
			expectedError: "file is empty or in incorrect format",
		},
		// Add more test cases like duplicate players, incorrect format, etc.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temp file and write the test content
			tmpFile, err := os.CreateTemp("", "testfile-*.txt")
			if err != nil {
				t.Fatal("Failed to create temp file:", err)
			}
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(tt.fileContent)
			if err != nil {
				t.Fatal("Failed to write to temp file:", err)
			}
			tmpFile.Close()

			// Now, call the function with the temp file path
			data, err := ReadInputFile(tmpFile.Name())

			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("Expected error %s, but got %s", tt.expectedError, err)
				return
			}

			if len(data) != len(tt.expectedData) {
				t.Errorf("Expected data length %d, but got %d", len(tt.expectedData), len(data))
				return
			}

			for i, row := range tt.expectedData {
				for j, player := range row {
					if data[i][j] != player {
						t.Errorf("Expected player %s, but got %s", player, data[i][j])
					}
				}
			}
		})
	}
}
