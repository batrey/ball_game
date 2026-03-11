package game

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuildGraph(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]map[string]bool
	}{
		{
			name:  "Basic input",
			input: "George,Beth,Sue\nRick,Anne\nAnne,Beth\nBeth,Anne,George\nSue,Beth\n",
			expected: map[string]map[string]bool{
				"George": {"Beth": true, "Sue": true},
				"Rick":   {"Anne": true},
				"Anne":   {"Beth": true},
				"Beth":   {"Anne": true, "George": true},
				"Sue":    {"Beth": true},
			},
		},
		{
			name:     "Empty input",
			input:    "",
			expected: map[string]map[string]bool{},
		},
		{
			name:  "Player can see self",
			input: "George,George,Beth,Sue\n",
			expected: map[string]map[string]bool{
				"George": {"George": true, "Beth": true, "Sue": true},
			},
		},
		{
			name:  "Multiple players, isolated",
			input: "A\nB\nC\n",
			expected: map[string]map[string]bool{
				"A": {},
				"B": {},
				"C": {},
			},
		},
		{
			name:  "Extra spaces",
			input: " George , Beth , Sue \nRick , Anne\n",
			expected: map[string]map[string]bool{
				"George": {"Beth": true, "Sue": true},
				"Rick":   {"Anne": true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := strings.Split(tt.input, "\n")
			var data [][]string
			for _, line := range lines {
				if line != "" {
					parts := strings.Split(line, ",")
					for i, part := range parts {
						parts[i] = strings.TrimSpace(part)
					}
					data = append(data, parts)
				}
			}

			graph, err := BuildGraph(data)
			if err != nil {
				t.Errorf("Expected no error, but got %s", err)
			}

			if !compareGraphs(graph, tt.expected) {
				t.Errorf("Expected graph to be %v but got %v", tt.expected, graph)
			}
		})
	}
}

func TestSpecialCharactersInNames(t *testing.T) {
	dataWithSpecialChars := [][]string{
		{"Al@n", "Br!an", "Car#ol"},
		{"Br&ian", "Al(an", "Diana%", "Ev^a"},
	}

	_, err := BuildGraph(dataWithSpecialChars)
	if err != nil {
		t.Fatalf("Expected no error for special characters in names, but got: %s", err)
	}
}

func compareGraphs(a, b map[string]map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if _, exists := b[k]; !exists {
			return false
		}
		for kk, vv := range v {
			if b[k][kk] != vv {
				return false
			}
		}
	}
	return true
}

func TestCalculateTouchesForPlayer(t *testing.T) {

	tests := []struct {
		name     string
		player   string
		graph    map[string]map[string]bool
		expected int
		want     int
	}{
		{
			name:   "George basic",
			player: "George",
			graph: map[string]map[string]bool{
				"George": {"Beth": true, "Sue": true},
				"Rick":   {"Anne": true},
				"Anne":   {"Beth": true},
				"Beth":   {"Anne": true, "George": true},
				"Sue":    {"Beth": true},
			},
			expected: 3,
		},
		{
			name:     "Isolated player",
			player:   "Sam",
			graph:    map[string]map[string]bool{"Sam": {}},
			expected: 1,
		},
		{
			name:     "Chain of players",
			player:   "A",
			graph:    map[string]map[string]bool{"A": {"B": true}, "B": {"C": true}, "C": {"D": true}, "D": {}},
			expected: 1, // Only 'A' can touch the ball since there's no mutual visibility between the players in the chain.
		},
		{
			name:     "Cycle of players",
			player:   "Adam",
			graph:    map[string]map[string]bool{"Adam": {"Eve": true, "Steve": true}, "Eve": {"Steve": true, "Adam": true}, "Steve": {"Adam": true, "Eve": true}},
			expected: 3,
		},
		{
			name:   "Star topology",
			player: "Center",
			graph: map[string]map[string]bool{
				"Center": {"A": true, "B": true, "C": true},
				"A":      {"Center": true},
				"B":      {"Center": true},
				"C":      {"Center": true},
			},
			expected: 4,
		},
		{
			name:   "Large graph, single connection",
			player: "Z",
			graph: func() map[string]map[string]bool {
				g := make(map[string]map[string]bool)
				for i := 'A'; i <= 'Z'; i++ {
					g[string(i)] = make(map[string]bool)
				}
				g["Z"]["Y"] = true
				g["Y"]["Z"] = true // Ensure mutual visibility
				return g
			}(),
			expected: 2,
		},
		{
			name:   "Graph with unconnected subgraphs",
			player: "A",
			graph: map[string]map[string]bool{
				"A": {"B": true},
				"B": {"A": true},
				"C": {"D": true},
				"D": {"C": true},
			},
			expected: 2,
		},
		{
			name:     "Player not present in the graph",
			player:   "E",
			graph:    map[string]map[string]bool{"A": {"B": true}, "B": {"A": true}},
			expected: 0,
		},
		{
			name:     "Visibility not reciprocal",
			player:   "A",
			graph:    map[string]map[string]bool{"A": {"B": true}, "B": {"A": true}},
			expected: 2,
		},
		{
			name:     "Graph with self visibility",
			player:   "A",
			graph:    map[string]map[string]bool{"A": {"A": true}},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := resetVisitedMap()

			result := CalculateTouchesForPlayer(tt.player, tt.graph)

			// Add debugging lines
			fmt.Printf("Test '%s': Visited map: %v\n", tt.name, visited)

			if result != tt.expected {
				t.Errorf("Expected %d players, but got %d", tt.expected, result)
			}
		})
	}
}

// A helper function to reset the visited map
func resetVisitedMap() map[string]bool {
	return make(map[string]bool)
}
