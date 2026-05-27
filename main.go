package main

import (
	"os"

	"ballgame/pkg/file"
	"ballgame/pkg/handler"
	"ballgame/pkg/service"
)

func main() {
	reader := file.Reader{}
	touchService := service.NewTouchService(reader)

	os.Exit(handler.RunCLI(os.Args[1:], os.Stdout, os.Stderr, touchService))
}
