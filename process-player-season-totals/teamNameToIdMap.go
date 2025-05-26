package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func teamNameIdMap() map[string]int {
	// API endpoint URL
	apiUrl := "https://statsapi.web.nhl.com/api/v1/teams?season=20212022"

	// Make an HTTP GET request to the API
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer response.Body.Close()

	// Parse the JSON response
	var data map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Fatalln("Error decoding JSON:", err)
	}

	// Check if the "teams" key exists in the JSON data
	teams, ok := data["teams"].([]interface{})
	if !ok {
		log.Fatalln("Error: 'teams' key not found in JSON data")
	}

	// Create a map to store teamName to teamId mappings
	teamNameToID := make(map[string]int)

	// Iterate through the teams and populate the map
	for _, team := range teams {
		if teamData, ok := team.(map[string]interface{}); ok {
			teamName := teamData["teamName"].(string)
			teamID := int(teamData["id"].(float64))
			teamNameToID[teamName] = teamID
		}
	}

	return teamNameToID
}
