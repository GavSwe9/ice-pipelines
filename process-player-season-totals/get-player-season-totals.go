package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gavswe19/ice-pipelines/database"
)

type SeasonTotals struct {
	Season             int      `json:"season"`
	GameTypeId         int      `json:"gameTypeId"`
	LeagueAbbrev       string   `json:"leagueAbbrev"`
	TeamName           string   `json:"teamName"`
	Sequence           int      `json:"sequence"`
	GamesPlayed        *int     `json:"gamesPlayed"`
	Goals              *int     `json:"goals"`
	Assists            *int     `json:"assists"`
	Points             *int     `json:"points"`
	PlusMinus          *int     `json:"plusMinus"`
	PowerPlayGoals     *int     `json:"powerPlayGoals"`
	PowerPlayPoints    *int     `json:"powerPlayPoints"`
	ShorthandedPoints  *int     `json:"shorthandedPoints"`
	GameWinningGoals   *int     `json:"gameWinningGoals"`
	OtGoals            *int     `json:"otGoals"`
	Shots              *int64   `json:"shots"`
	ShootingPctg       *float64 `json:"shootingPctg"`
	FaceoffWinningPctg *float64 `json:"faceoffWinningPctg"`
	AvgToi             *string  `json:"avgToi"`
	ShorthandedGoals   *int     `json:"shorthandedGoals"`
	Pim                *int     `json:"pim"`
}

type PlayerData struct {
	SeasonTotals []SeasonTotals `json:"seasonTotals"`
}

func main() {
	lambda.Start(Handler)
}

// func main() {
// 	var ctx context.Context
// 	var sqsEvent events.SQSEvent

// 	message := events.SQSMessage{
// 		MessageId: "mock-message-id",
// 		Body:      "8478402",
// 	}

// 	sqsEvent.Records = []events.SQSMessage{message}

// 	Handler(ctx, sqsEvent)
// }

func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	eventRecord := sqsEvent.Records[0]
	body := eventRecord.Body
	playerId, err := strconv.Atoi(body)

	db := database.GetDatabase()
	tx, err := db.Begin()
	println("Start Transaction")

	processPlayerSeasonTotals(tx, playerId)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	println("Committed Transaction")
}

func processPlayerSeasonTotals(tx *sql.Tx, playerId int) {
	teamNameToID := teamNameIdMap()

	// API endpoint URL
	apiUrl := fmt.Sprintf("https://api-web.nhle.com/v1/player/%d/landing", playerId)

	// Make an HTTP GET request to the API
	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create a struct to hold the parsed JSON data
	var playerData PlayerData

	// Parse the JSON response into the struct
	err = json.Unmarshal(responseBody, &playerData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Access the list of seasonTotals
	seasonTotals := playerData.SeasonTotals

	insertPlayerSeasonTotals(tx, playerId, seasonTotals, teamNameToID)
}

func insertPlayerSeasonTotals(tx *sql.Tx, playerId int, seasonTotals []SeasonTotals, teamNameToID map[string]int) {
	playerSeasonTotalsStrings := make([]string, 0, len(seasonTotals))
	playerSeasonTotalsValueArgs := make([]interface{}, 0, len(seasonTotals)*9)

	for _, line := range seasonTotals {
		playerSeasonTotalsStrings = append(playerSeasonTotalsStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		playerSeasonTotalsValueArgs = append(playerSeasonTotalsValueArgs,
			playerId,
			line.Season,
			teamNameToID[line.TeamName],
			line.GameTypeId,
			line.LeagueAbbrev,
			line.TeamName,
			line.Sequence,
			line.GamesPlayed,
			line.Shots,
			line.Goals,
			line.Assists,
			line.Points,
			line.PlusMinus,
			line.PowerPlayGoals,
			line.PowerPlayPoints,
			line.ShorthandedGoals,
			line.ShorthandedPoints,
			line.GameWinningGoals,
			line.OtGoals,
			line.ShootingPctg,
			line.FaceoffWinningPctg,
			line.AvgToi,
			line.Pim,
		)
		fmt.Println(line.LeagueAbbrev, line.Shots)
	}

	stmt := fmt.Sprintf("INSERT INTO player_season_totals VALUES %s ON DUPLICATE KEY UPDATE player_id=player_id", strings.Join(playerSeasonTotalsStrings, ","))
	result, err := tx.Exec(stmt, playerSeasonTotalsValueArgs...)

	if err != nil {
		tx.Rollback()
		fmt.Println("Rolling back")
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into player_season_totals", strconv.Itoa(int(rows_affected))))
}
