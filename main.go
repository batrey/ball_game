package main

import (
	"flag"
	"fmt"
	"log"

	"ballgame/pkg/game"

	"ballgame/pkg/file"
)

func main() {
	fileName := flag.String("file", "", "Path to the input file")
	flag.Parse()

	if *fileName == "" {
		log.Fatal("Please provide a valid input file using the -file flag.")
		return
	}

	playersData, err := file.ReadInputFile(*fileName)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
		return
	}

	adjList, err := game.BuildGraph(playersData)
	if err != nil {
		log.Fatalf("Error building graph: %s", err)
		return
	}

	// Using goroutines for parallel processing (concurrency)
	resultChan := make(chan int)
	for player, visiblePlayers := range adjList {
		go func(p string, v map[string]bool) {
			touchCount := game.CalculateTouchesForPlayer(p, adjList)
			resultChan <- touchCount
		}(player, visiblePlayers)
	}

	maxTouch := 0
	for range adjList {
		touchCount := <-resultChan
		if touchCount > maxTouch {
			maxTouch = touchCount
		}
	}

	fmt.Println(maxTouch)
}
