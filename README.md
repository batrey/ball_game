# Ball Game Touch Counter

Ball Game Touch Counter reads a comma-separated visibility file and prints the maximum number of players who can touch the ball. A player can pass the ball only through mutual visibility: if `A` lists `B`, `B` must also list `A`.

## Input Format

Each non-blank line starts with a player, followed by the players they can see:

```text
George,Beth,Sue
Rick,Anne
Anne,Beth
Beth,Anne,George
Sue,Beth
```

Rules:

- Names are trimmed for surrounding whitespace.
- Blank lines are ignored.
- Empty names are invalid.
- Names may contain special or Unicode characters.
- Names containing commas must use CSV quoting, such as `"George, Jr.",Beth`.
- A name can be at most 20 characters.

## Run

Requires Go 1.19 or newer.

```bash
go run . -file /path/to/players.txt
```

Build a binary:

```bash
go build -o ballgame
./ballgame -file /path/to/players.txt
```

## Test

```bash
go test ./...
```

## Project Structure

```text
main.go              wires the CLI handler, service, and file reader
pkg/handler          command-line boundary: flags, stdout/stderr, exit codes
pkg/service          orchestration: read player data, build graph, calculate result
pkg/file             file-backed data access
pkg/game             domain logic: graph building and mutual-visibility traversal
```

## Algorithm

The app builds an adjacency list from the input file, then finds the largest connected component using only reciprocal edges. For example, `A -> B` counts only when `B -> A` also exists.

`MaxTouches` visits each graph component once, so the calculation is linear in the number of players plus visibility entries.
