package game

import (
	"errors"
	"log"
)

func BuildGraph(playersData [][]string) (map[string]map[string]bool, error) {
	adjList := make(map[string]map[string]bool)

	for _, row := range playersData {
		player := row[0]

		// Check for maximum length
		if len(player) > 20 {
			return nil, errors.New("player name exceeds maximum length of 20 characters")
		}

		if _, exists := adjList[player]; !exists {
			adjList[player] = make(map[string]bool)
		}
		for _, p := range row[1:] {
			// Check for maximum length
			if len(p) > 20 {
				return nil, errors.New("player name exceeds maximum length of 20 characters")
			}
			adjList[player][p] = true
		}
	}
	return adjList, nil
}

// CalculateTouchesForPlayer starts from a player and calculates how many players can touch the ball, considering mutual visibility.
func CalculateTouchesForPlayer(player string, graph map[string]map[string]bool) int {
	visited := make(map[string]bool)
	log.Printf("Starting calculation for player: %s\n", player)
	count := dfs(player, graph, visited)
	log.Printf("Total players that can touch the ball (starting from %s): %d\n", player, count)
	return count
}

// dfs is a depth-first search function that traverses the graph based on mutual visibility and returns the number of players that can touch the ball.
func dfs(player string, graph map[string]map[string]bool, visited map[string]bool) int {
	if _, exists := graph[player]; !exists { // Player is not in the graph
		log.Printf("Player %s not present in the graph.\n", player)
		return 0
	}

	if visited[player] { // Player has already been visited
		log.Printf("Player %s has already been visited.\n", player)
		return 0
	}

	visited[player] = true
	count := 1 // Current player can touch the ball
	log.Printf("Player %s can touch the ball. Exploring visible players...\n", player)

	for adjPlayer := range graph[player] {
		if graph[adjPlayer][player] { // Ensure mutual visibility
			log.Printf("Mutual visibility found between %s and %s.\n", player, adjPlayer)
			count += dfs(adjPlayer, graph, visited)
		} else {
			log.Printf("%s can see %s, but the visibility is not mutual.\n", player, adjPlayer)
		}
	}

	return count
}
