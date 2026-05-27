package game

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const maxPlayerNameLength = 20

// BuildGraph builds an adjacency list from player visibility rows.
func BuildGraph(playersData [][]string) (map[string]map[string]bool, error) {
	adjList := make(map[string]map[string]bool)

	for _, row := range playersData {
		if len(row) == 0 {
			return nil, errors.New("row is empty or in incorrect format")
		}

		player, err := normalizePlayerName(row[0])
		if err != nil {
			return nil, err
		}

		if _, exists := adjList[player]; !exists {
			adjList[player] = make(map[string]bool)
		}
		for _, p := range row[1:] {
			visiblePlayer, err := normalizePlayerName(p)
			if err != nil {
				return nil, err
			}
			adjList[player][visiblePlayer] = true
		}
	}
	return adjList, nil
}

// CalculateTouchesForPlayer returns how many players can touch the ball from the given starting player.
func CalculateTouchesForPlayer(player string, graph map[string]map[string]bool) int {
	visited := make(map[string]bool)
	return countReachablePlayers(player, graph, visited)
}

// MaxTouches returns the largest mutual-visibility component in the graph.
func MaxTouches(graph map[string]map[string]bool) int {
	visited := make(map[string]bool)
	maxTouches := 0

	for player := range graph {
		if visited[player] {
			continue
		}

		touches := countReachablePlayers(player, graph, visited)
		if touches > maxTouches {
			maxTouches = touches
		}
	}

	return maxTouches
}

func countReachablePlayers(player string, graph map[string]map[string]bool, visited map[string]bool) int {
	if _, exists := graph[player]; !exists || visited[player] {
		return 0
	}

	visited[player] = true
	count := 1

	for adjPlayer := range graph[player] {
		if canSeeEachOther(player, adjPlayer, graph) {
			count += countReachablePlayers(adjPlayer, graph, visited)
		}
	}

	return count
}

func canSeeEachOther(first string, second string, graph map[string]map[string]bool) bool {
	return graph[first][second] && graph[second][first]
}

func normalizePlayerName(name string) (string, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return "", errors.New("player name cannot be empty")
	}
	if utf8.RuneCountInString(trimmedName) > maxPlayerNameLength {
		return "", fmt.Errorf("player name exceeds maximum length of %d characters", maxPlayerNameLength)
	}
	return trimmedName, nil
}
