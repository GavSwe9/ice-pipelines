package main

import (
	"fmt"
	"log"
)

func main() {
	// Initialize repository using your existing database connection
	repo, err := NewPlayerRepository()
	if err != nil {
		log.Fatalf("Failed to initialize database repository: %v", err)
	}

	// Create table if it doesn't exist
	fmt.Println("Creating table if not exists...")
	if err := repo.CreateTable(); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table ready.")

	// Parse skaters CSV
	fmt.Println("Parsing skaters CSV...")
	skaters, err := parseSkaterCSV("csv/skaters_07to25.csv")
	if err != nil {
		log.Fatalf("Failed to parse skaters CSV: %v", err)
	}
	fmt.Printf("Parsed %d skater records.\n", len(skaters))

	// Parse goalies CSV
	fmt.Println("Parsing goalies CSV...")
	goalies, err := parseGoalieCSV("csv/goalies_07to25.csv")
	if err != nil {
		log.Fatalf("Failed to parse goalies CSV: %v", err)
	}
	fmt.Printf("Parsed %d goalie records.\n", len(goalies))

	// Combine all players
	allPlayers := make([]PlayerStats, 0, len(skaters)+len(goalies))
	allPlayers = append(allPlayers, skaters...)
	allPlayers = append(allPlayers, goalies...)
	fmt.Printf("Total records to insert: %d\n", len(allPlayers))

	// Insert into database
	fmt.Println("Inserting records into database...")
	if err := repo.InsertPlayerStats(allPlayers); err != nil {
		log.Fatalf("Failed to insert player stats: %v", err)
	}

	fmt.Println("Data processing completed successfully!")
}
