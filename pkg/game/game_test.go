package game

import (
	"reflect"
	"strings"
	"testing"
)

func TestBuildGraph(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]string
		expected map[string]map[string]bool
	}{
		{
			name: "basic input",
			input: [][]string{
				{"George", "Beth", "Sue"},
				{"Rick", "Anne"},
				{"Anne", "Beth"},
				{"Beth", "Anne", "George"},
				{"Sue", "Beth"},
			},
			expected: map[string]map[string]bool{
				"George": {"Beth": true, "Sue": true},
				"Rick":   {"Anne": true},
				"Anne":   {"Beth": true},
				"Beth":   {"Anne": true, "George": true},
				"Sue":    {"Beth": true},
			},
		},
		{
			name:     "empty input",
			input:    nil,
			expected: map[string]map[string]bool{},
		},
		{
			name:  "player can see self",
			input: [][]string{{"George", "George", "Beth", "Sue"}},
			expected: map[string]map[string]bool{
				"George": {"George": true, "Beth": true, "Sue": true},
			},
		},
		{
			name: "multiple isolated players",
			input: [][]string{
				{"A"},
				{"B"},
				{"C"},
			},
			expected: map[string]map[string]bool{
				"A": {},
				"B": {},
				"C": {},
			},
		},
		{
			name: "trims extra spaces",
			input: [][]string{
				{" George ", " Beth ", " Sue "},
				{"Rick ", " Anne"},
			},
			expected: map[string]map[string]bool{
				"George": {"Beth": true, "Sue": true},
				"Rick":   {"Anne": true},
			},
		},
		{
			name: "deduplicates repeated visible players",
			input: [][]string{
				{"George", "Beth", "Beth"},
			},
			expected: map[string]map[string]bool{
				"George": {"Beth": true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given / when
			graph, err := BuildGraph(tt.input)

			// then
			if err != nil {
				t.Fatalf("expected no error, got %s", err)
			}
			if !reflect.DeepEqual(graph, tt.expected) {
				t.Fatalf("expected graph %v, got %v", tt.expected, graph)
			}
		})
	}
}

func TestBuildGraphRejectsInvalidRows(t *testing.T) {
	tests := []struct {
		name        string
		input       [][]string
		expectedErr string
	}{
		{
			name:        "empty row",
			input:       [][]string{{}},
			expectedErr: "row is empty or in incorrect format",
		},
		{
			name:        "empty player name",
			input:       [][]string{{" "}},
			expectedErr: "player name cannot be empty",
		},
		{
			name:        "empty visible player name",
			input:       [][]string{{"Alice", " "}},
			expectedErr: "player name cannot be empty",
		},
		{
			name:        "too long player name",
			input:       [][]string{{strings.Repeat("A", 21)}},
			expectedErr: "player name exceeds maximum length of 20 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given / when
			_, err := BuildGraph(tt.input)

			// then
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if err.Error() != tt.expectedErr {
				t.Fatalf("expected error %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestBuildGraphCountsUnicodeNameLengthByCharacter(t *testing.T) {
	// given
	playerName := strings.Repeat("å", 20)

	// when
	graph, err := BuildGraph([][]string{{playerName}})

	// then
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if _, exists := graph[playerName]; !exists {
		t.Fatalf("expected graph to contain %q, got %v", playerName, graph)
	}
}

func TestCalculateTouchesForPlayer(t *testing.T) {
	tests := []struct {
		name     string
		player   string
		graph    map[string]map[string]bool
		expected int
	}{
		{
			name:   "basic mutual visibility",
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
			name:     "isolated player",
			player:   "Sam",
			graph:    map[string]map[string]bool{"Sam": {}},
			expected: 1,
		},
		{
			name:     "chain without reciprocal visibility",
			player:   "A",
			graph:    map[string]map[string]bool{"A": {"B": true}, "B": {"C": true}, "C": {"D": true}, "D": {}},
			expected: 1,
		},
		{
			name:     "cycle of players",
			player:   "Adam",
			graph:    map[string]map[string]bool{"Adam": {"Eve": true, "Steve": true}, "Eve": {"Steve": true, "Adam": true}, "Steve": {"Adam": true, "Eve": true}},
			expected: 3,
		},
		{
			name:   "star topology",
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
			name:     "player not present in graph",
			player:   "E",
			graph:    map[string]map[string]bool{"A": {"B": true}, "B": {"A": true}},
			expected: 0,
		},
		{
			name:     "direct visibility must be reciprocal",
			player:   "A",
			graph:    map[string]map[string]bool{"A": {"B": true}, "B": {}},
			expected: 1,
		},
		{
			name:     "self visibility does not double count",
			player:   "A",
			graph:    map[string]map[string]bool{"A": {"A": true}},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given / when
			result := CalculateTouchesForPlayer(tt.player, tt.graph)

			// then
			if result != tt.expected {
				t.Fatalf("expected %d players, got %d", tt.expected, result)
			}
		})
	}
}

func TestMaxTouches(t *testing.T) {
	tests := []struct {
		name     string
		graph    map[string]map[string]bool
		expected int
	}{
		{
			name:     "empty graph",
			graph:    map[string]map[string]bool{},
			expected: 0,
		},
		{
			name: "returns largest mutual component",
			graph: map[string]map[string]bool{
				"A": {"B": true},
				"B": {"A": true},
				"C": {"D": true},
				"D": {"C": true, "E": true},
				"E": {"D": true},
			},
			expected: 3,
		},
		{
			name: "ignores one-way visibility",
			graph: map[string]map[string]bool{
				"A": {"B": true},
				"B": {},
				"C": {"D": true},
				"D": {"C": true},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given / when
			result := MaxTouches(tt.graph)

			// then
			if result != tt.expected {
				t.Fatalf("expected %d players, got %d", tt.expected, result)
			}
		})
	}
}
