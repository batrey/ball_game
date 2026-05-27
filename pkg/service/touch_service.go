package service

import "ballgame/pkg/game"

// PlayerDataReader loads player visibility rows from storage.
type PlayerDataReader interface {
	ReadPlayers(fileName string) ([][]string, error)
}

// TouchService coordinates player data loading and touch-count calculation.
type TouchService struct {
	reader PlayerDataReader
}

// NewTouchService creates a service backed by the given player data reader.
func NewTouchService(reader PlayerDataReader) TouchService {
	return TouchService{reader: reader}
}

// MaxTouchesFromFile returns the maximum touch count for player data in a file.
func (s TouchService) MaxTouchesFromFile(fileName string) (int, error) {
	playersData, err := s.reader.ReadPlayers(fileName)
	if err != nil {
		return 0, err
	}

	graph, err := game.BuildGraph(playersData)
	if err != nil {
		return 0, err
	}

	return game.MaxTouches(graph), nil
}
