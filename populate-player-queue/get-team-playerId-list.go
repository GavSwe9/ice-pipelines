package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Define the structures that represent the JSON response

type PlayerStatsResponse struct {
	Season   string   `json:"season"`
	GameType int      `json:"gameType"`
	Skaters  []Player `json:"skaters"`
	Goalies  []Player `json:"goalies"`
}

type Player struct {
	PlayerId int `json:"playerId"`
}

func getTeamPlayerIdList(teamAbbrev string) []int {
	url := fmt.Sprintf("https://api-web.nhle.com/v1/club-stats/%s/20212022/2", teamAbbrev)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data from API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var stats PlayerStatsResponse
	err = json.Unmarshal(body, &stats)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	var playerIds []int
	for _, skater := range stats.Skaters {
		playerIds = append(playerIds, skater.PlayerId)
	}
	for _, goalie := range stats.Goalies {
		playerIds = append(playerIds, goalie.PlayerId)
	}

	fmt.Println("Player IDs:")
	for _, id := range playerIds {
		fmt.Println(id)
	}

	return playerIds
}
