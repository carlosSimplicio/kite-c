package main

import (
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"kite-c/components"
)

func main() {
	logOutputFile := "logs.txt"
	file, err := os.OpenFile(logOutputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	log.SetOutput(file)

	p := tea.NewProgram(components.NewWindowModel())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Failed to run program %s", err.Error())
	}
}
