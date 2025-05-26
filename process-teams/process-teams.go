package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gavswe19/ice-pipelines/database"
)

type TeamInfo struct {
	Default string `json:"default"`
	Fr      string `json:"fr,omitempty"`
}

type Standings struct {
	SeasonId       int      `json:"seasonId"`
	ConferenceName string   `json:"conferenceName"`
	DivisionName   string   `json:"divisionName"`
	TeamName       TeamInfo `json:"teamName"`
	TeamCommonName TeamInfo `json:"teamCommonName"`
	TeamAbbrev     TeamInfo `json:"teamAbbrev"`
}

type StandingsResponse struct {
	WildCardIndicator bool        `json:"wildCardIndicator"`
	Standings         []Standings `json:"standings"`
}

func main() {
	db := database.GetDatabase()
	url := "https://api-web.nhle.com/v1/standings/2024-01-01"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var standingsResponse StandingsResponse
	err = json.Unmarshal(body, &standingsResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	for _, standing := range standingsResponse.Standings {
		fmt.Printf("%s (%s)\n", standing.TeamName.Default, standing.TeamAbbrev.Default)

		var id int
		err = db.QueryRow("SELECT id FROM team_seasons_2 WHERE team_abbrev = ? AND season_id = ?", standing.TeamAbbrev.Default, standing.SeasonId).Scan(&id)

		if err == sql.ErrNoRows {
			// Record does not exist, insert a new one
			_, err = db.Exec(
				"INSERT INTO team_seasons_2 (season_id, conference_name, division_name, team_name, team_common_name, team_abbrev) VALUES (?, ?, ?, ?, ?, ?)",
				standing.SeasonId, standing.ConferenceName, standing.DivisionName, standing.TeamName.Default, standing.TeamCommonName.Default, standing.TeamAbbrev.Default,
			)
			if err != nil {
				fmt.Println("Error inserting new record:", err)
				continue
			}
			fmt.Printf("Inserted new record for team %s season %d\n", standing.TeamName.Default, standing.SeasonId)
		} else if err != nil {
			fmt.Println("Error querying for existing record:", err)
		}
	}
}
