package service

import (
	"errors"
	"testing"
)

func TestTouchServiceMaxTouchesFromFile(t *testing.T) {
	// given
	reader := &mockPlayerDataReader{
		data: [][]string{
			{"A", "B"},
			{"B", "A"},
			{"C", "D"},
			{"D", "C", "E"},
			{"E", "D"},
		},
	}
	service := NewTouchService(reader)

	// when
	result, err := service.MaxTouchesFromFile("players.txt")

	// then
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if result != 3 {
		t.Fatalf("expected 3 touches, got %d", result)
	}
	if reader.gotFileName != "players.txt" {
		t.Fatalf("expected reader to receive file name %q, got %q", "players.txt", reader.gotFileName)
	}
}

func TestTouchServiceReturnsReaderError(t *testing.T) {
	// given
	expectedErr := errors.New("open failed")
	reader := &mockPlayerDataReader{err: expectedErr}
	service := NewTouchService(reader)

	// when
	_, err := service.MaxTouchesFromFile("missing.txt")

	// then
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %q, got %v", expectedErr, err)
	}
}

func TestTouchServiceReturnsBuildGraphError(t *testing.T) {
	// given
	reader := &mockPlayerDataReader{data: [][]string{{"Alice", ""}}}
	service := NewTouchService(reader)

	// when
	_, err := service.MaxTouchesFromFile("players.txt")

	// then
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != "player name cannot be empty" {
		t.Fatalf("expected validation error, got %q", err.Error())
	}
}

type mockPlayerDataReader struct {
	data        [][]string
	err         error
	gotFileName string
}

func (m *mockPlayerDataReader) ReadPlayers(fileName string) ([][]string, error) {
	m.gotFileName = fileName
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}
