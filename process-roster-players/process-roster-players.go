package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func processTeamRoster(tx *sql.Tx, season string) TeamsResponse {
	response, err := http.Get(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/teams?season=%s", season))

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var teamsResponse TeamsResponse
	json.Unmarshal(responseData, &teamsResponse)

	insertTeamSeasons(tx, teamsResponse.Teams, season)
	for _, team := range teamsResponse.Teams {
		rosterPlayers := getRosterPlayers(team.Id, season)
		insertPlayers(tx, rosterPlayers)
		insertTeamSeasonPlayers(tx, rosterPlayers, team.Id, season)
	}

	return teamsResponse
}

func getRosterPlayers(teamId int, season string) []PlayerObj {
	response, err := http.Get(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/teams/%s/roster?season=%s", strconv.Itoa(teamId), season))

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var rosterResponse RosterResponse
	json.Unmarshal(responseData, &rosterResponse)
	return rosterResponse.Players
}

func insertPlayers(tx *sql.Tx, players []PlayerObj) {
	playerStrings := make([]string, 0, len(players))
	playerValueArgs := make([]interface{}, 0, len(players)*3)

	for _, player := range players {
		playerStrings = append(playerStrings, "(?, ?, ?)")
		playerValueArgs = append(playerValueArgs,
			player.Player.PlayerId,
			player.Player.FullName,
			player.Position.PositionCode,
		)
	}

	stmt := fmt.Sprintf("INSERT INTO players VALUES %s ON DUPLICATE KEY UPDATE position=position", strings.Join(playerStrings, ","))
	result, err := tx.Exec(stmt, playerValueArgs...)

	if err != nil {
		tx.Rollback()
		fmt.Println("Rolling back")
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into players", strconv.Itoa(int(rows_affected))))
}

func insertTeamSeasonPlayers(tx *sql.Tx, players []PlayerObj, teamId int, season string) {
	teamSeasonPlayerStrings := make([]string, 0, len(players))
	teamSeasonPlayerValueArgs := make([]interface{}, 0, len(players)*3)

	for _, player := range players {
		teamSeasonPlayerStrings = append(teamSeasonPlayerStrings, "(?, ?, ?)")
		teamSeasonPlayerValueArgs = append(teamSeasonPlayerValueArgs,
			teamId,
			season,
			player.Player.PlayerId,
		)
	}

	stmt := fmt.Sprintf("INSERT INTO team_season_players VALUES %s ON DUPLICATE KEY UPDATE team_id=team_id", strings.Join(teamSeasonPlayerStrings, ","))
	result, err := tx.Exec(stmt, teamSeasonPlayerValueArgs...)

	if err != nil {
		tx.Rollback()
		fmt.Println("Rolling back")
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into team_season_players", strconv.Itoa(int(rows_affected))))
}

func insertTeamSeasons(tx *sql.Tx, teams []Team, season string) {
	teamSeasonsStrings := make([]string, 0, len(teams))
	teamSeasonsValueArgs := make([]interface{}, 0, len(teams)*9)

	for _, team := range teams {
		teamSeasonsStrings = append(teamSeasonsStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		teamSeasonsValueArgs = append(teamSeasonsValueArgs,
			season,
			team.Id,
			team.Name,
			team.Abbreviation,
			team.Division.Id,
			team.Division.Name,
			team.Conference.Id,
			team.Conference.Name,
			team.FranchiseId,
		)
	}

	stmt := fmt.Sprintf("INSERT INTO team_seasons VALUES %s ON DUPLICATE KEY UPDATE team_id=team_id", strings.Join(teamSeasonsStrings, ","))
	result, err := tx.Exec(stmt, teamSeasonsValueArgs...)

	if err != nil {
		tx.Rollback()
		fmt.Println("Rolling back")
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into team_seasons", strconv.Itoa(int(rows_affected))))
}

func main() {
	secrets := getAwsSecrets()

	cfg := mysql.Config{
		User:                 secrets.Username,
		Passwd:               secrets.Password,
		Net:                  "tcp",
		Addr:                 "farm.cxqsjcdo8n1w.us-east-1.rds.amazonaws.com",
		DBName:               "ICE",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	println("Start Transation")

	season := "20202021"
	processTeamRoster(tx, season)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	println("Committed Transation")
}
