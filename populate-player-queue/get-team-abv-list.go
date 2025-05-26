package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Define the structures that represent the JSON response

type StandingsResponse struct {
	Teams []Team `json:"standings"`
}

type Team struct {
	Abbreviation Abbreviation `json:"teamAbbrev"`
}

type Abbreviation struct {
	Default string `json:"default"`
}

func getAllTeamAbv() []string {
	// URL of the API
	url := "https://api-web.nhle.com/v1/standings/2024-01-01"

	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data from API: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Parse the JSON response
	var standings StandingsResponse
	err = json.Unmarshal(body, &standings)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Extract the list of teamAbbrev.default values
	var teamAbbrevs []string
	for _, team := range standings.Teams {
		teamAbbrevs = append(teamAbbrevs, team.Abbreviation.Default)
	}

	return teamAbbrevs
}
