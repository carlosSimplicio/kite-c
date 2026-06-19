package main

import (
	tea "charm.land/bubbletea/v2"
	"context"
	"database/sql"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
	"kite-c/components"
	commandRepository "kite-c/services/repository"
	"log"
	"os"
)

//go:embed database/schema.sql
var ddl string

func main() {
	ctx := context.Background()
	dbClient, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer dbClient.Close()

	if err := dbClient.PingContext(ctx); err != nil {
		log.Fatalf("failed to reach database: %v", err)
	}

	if _, err := dbClient.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("failed to enable sqlite foreign keys: %v", err)
	}

	if _, err := dbClient.ExecContext(ctx, ddl); err != nil {
		log.Fatalln("Failed to migrate", err.Error())
	}

	repository := commandRepository.NewCommandRepository(dbClient)

	logOutputFile := "logs.txt"
	file, err := os.OpenFile(logOutputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	log.SetOutput(file)

	p := tea.NewProgram(components.NewWindowModel(repository))

	if _, err := p.Run(); err != nil {
		log.Fatalf("Failed to run program %s", err.Error())
	}
}
