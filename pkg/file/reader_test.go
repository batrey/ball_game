package file

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadInputFile(t *testing.T) {
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
		},
		{
			name: "trims spaces and skips blank lines",
			fileContent: `
 George , Beth , Sue

 Rick , Anne
`,
			expectedData: [][]string{
				{"George", "Beth", "Sue"},
				{"Rick", "Anne"},
			},
		},
		{
			name:          "empty file",
			fileContent:   "",
			expectedData:  nil,
			expectedError: "file is empty or in incorrect format",
		},
		{
			name:          "whitespace only file",
			fileContent:   " \n\t\n",
			expectedData:  nil,
			expectedError: "file is empty or in incorrect format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			fileName := writeTempFile(t, tt.fileContent)

			// when
			data, err := ReadInputFile(fileName)

			// then
			if tt.expectedError != "" {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				if err.Error() != tt.expectedError {
					t.Fatalf("expected error %q, got %q", tt.expectedError, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %s", err)
			}
			if !reflect.DeepEqual(data, tt.expectedData) {
				t.Fatalf("expected data %v, got %v", tt.expectedData, data)
			}
		})
	}
}

func TestReadInputFileReturnsOpenError(t *testing.T) {
	// given
	fileName := filepath.Join(t.TempDir(), "missing.txt")

	// when
	_, err := ReadInputFile(fileName)

	// then
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestReaderReadPlayers(t *testing.T) {
	// given
	fileName := writeTempFile(t, "George,Beth\n")
	reader := Reader{}

	// when
	data, err := reader.ReadPlayers(fileName)

	// then
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	expected := [][]string{{"George", "Beth"}}
	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("expected data %v, got %v", expected, data)
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	fileName := filepath.Join(t.TempDir(), "players.txt")
	if err := os.WriteFile(fileName, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %s", err)
	}

	return fileName
}
