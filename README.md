# Ball Game: Touch Count Calculator

This application calculates the maximum number of players who can touch a single ball based on mutual visibility. It takes an input file (txt) where each line denotes a player and the players they can see. It uses graph theory principles and a depth-first search algorithm to determine the touch count for each player.

## Time Complexity

The expected time complexity of the algorithm for a file with \( n \) players, where each player can see at most \( m \) other players, is \( O(n \times m) \).

In the provided implementation:

- Parsing the file has a time complexity of \( O(n) \) since we read the file line by line.
- Building the adjacency list (graph) is \( O(n \times m) \) as we iterate over each player and then the players they can see.
- The depth-first search (DFS) also runs at \( O(n \times m) \) as, in the worst case, we might end up visiting all players and all their adjacent players.

Hence, combining all steps, our solution's time complexity is \( O(n \times m) \).

## Error Handling

The application is equipped to handle erroneous input data:

- If a player's name exceeds 20 characters, a relevant error message is thrown.
- If the input file is not in the expected format or is empty, the program will notify the user with a descriptive error message.

## Language and Execution

The solution is written in Go (Golang).

### Prerequisites

- Ensure you have Go 1.19 installed. If not, download and install it from [here](https://golang.org/dl/).

### Compiling and Running

1. Navigate to the directory containing the code.
2. Run the application using:

```bash
go run main.go -file /path/to/your/inputfile.txt
```

or

Build the Game:

```bash
go build -o ballgame
```

This will produce an executable named ballgame in your project directory.

Running the Game
After building the project, you can run the game using the following command:

On Windows, you might run:

```bash
ballgame.exe
```

```bash
go run ballgame -file /path/to/your/input
```

Replace `/path/to/your/inputfile.txt` with the path to your input file.

It looks like you have comprehensive test suites for different components of your program. Here's a rewritten testing section for your README based on the tests you provided:

---

## Testing

### How to Run the Tests

To run the provided tests, navigate to the root directory of the project and execute:

```bash
go test ./...
```

### Test Overview

We have tests for both the core game logic and the file reader.

#### Game Tests (`game_test.go`)

1. **Graph Building (`TestBuildGraph`)**: This test verifies that the `BuildGraph` function processes player visibility correctly and forms the right graph structure. Test cases include:
   - Basic input
   - Empty input
   - Players being able to see themselves
   - Multiple isolated players
   - Extra spaces in input

2. **Special Character Handling (`TestSpecialCharactersInNames`)**: Tests the ability of the system to handle special characters in player names without errors.

3. **Touch Calculation (`TestCalculateTouchesForPlayer`)**:
   - This test checks the logic for counting how many players a given player can touch.
   - Scenarios tested include basic visibility, isolated players, chains of players, cycles in visibility, star topology, large graphs, unconnected subgraphs, players not present in the graph, non-reciprocal visibility, and self visibility.

#### File Reader Tests (`reader_test.go`)

1. **File Reading (`TestReadInputFile`)**:
   - This test verifies that the file reading function can correctly read and parse different formats of input files.
   - Test cases include basic input and empty files. You can expand this test with more cases like duplicate players, incorrect file format, etc.
