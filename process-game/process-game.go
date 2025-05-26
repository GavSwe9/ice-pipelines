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

// var db *sql.DB

func main() {
	lambda.Start(Handler)
}

// func main() {
// 	var ctx context.Context
// 	var sqsEvent events.SQSEvent

// 	message := events.SQSMessage{
// 		MessageId: "mock-message-id",
// 		Body:      "2022020287",
// 	}

// 	sqsEvent.Records = []events.SQSMessage{message}

// 	Handler(ctx, sqsEvent)
// }

func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	eventRecord := sqsEvent.Records[0]

	body := eventRecord.Body

	gamePk, err := strconv.Atoi(body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf(" *** Processing GamePk %s ***", strconv.Itoa(gamePk)))

	db := database.GetDatabase()
	gameAlreadyProcessed := gameHasBeenProcessed(db, gamePk)
	if gameAlreadyProcessed {
		fmt.Println(fmt.Sprintf("GamePk %s has already been processed", strconv.Itoa(gamePk)))
		return
	}

	UpdateEtlGameStatus(db, gamePk, "IN_PROGRESS")

	response, err := http.Get(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/game/%s/feed/live", strconv.Itoa(gamePk)))

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responseObject GameResponse
	json.Unmarshal(responseData, &responseObject)

	onIceRecordList := GetPlayersOnIce(gamePk, responseObject.LiveData.Plays.AllPlays)

	tx, err := db.Begin()
	println("Start Transation")

	// DeleteWithGamePk(tx, "play_by_play", gamePk)
	// DeleteWithGamePk(tx, "play_by_play_contributor", gamePk)
	// DeleteWithGamePk(tx, "play_by_play_on_ice", gamePk)

	season, err := strconv.Atoi(responseObject.GameData.Game.Season)
	if err != nil {
		log.Fatal("Error parsing season string to int")
	}

	DeleteWithGamePk(tx, "games", gamePk)
	InsertGames(tx, responseObject.GameData)
	InsertPlayByPlayRecords(tx, gamePk, responseObject.LiveData.Plays.AllPlays)
	InsertOnIceRecords(tx, onIceRecordList)
	InsertSkaterLineRecords(db, onIceRecordList, season)

	UpdateEtlGameStatus(db, gamePk, "COMPLETE")

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	println("Committed Transation")

}

func DeleteWithGamePk(tx *sql.Tx, tableName string, gamePk int) {
	result, err := tx.Exec(fmt.Sprintf("DELETE FROM %s PBP WHERE PBP.game_pk = %s", tableName, strconv.Itoa(gamePk)))
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Delted %s records from %s for gamePk %s", strconv.Itoa(int(rows_affected)), tableName, strconv.Itoa(gamePk)))
}

func InsertGames(tx *sql.Tx, game GameData) {
	stmt := fmt.Sprintf("INSERT INTO games VALUES (\"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\")",
		strconv.Itoa(game.Game.GamePk),
		game.Game.Type,
		game.Game.Season,
		game.DateTime.DateTime,
		strconv.Itoa(game.Teams.AwayTeam.Id),
		strconv.Itoa(game.Teams.HomeTeam.Id),
	)

	println(stmt)
	result, err := tx.Exec(stmt)

	if err != nil {
		tx.Rollback()
		fmt.Println("Rolling back")
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s record into games for gamePk %s", strconv.Itoa(int(rows_affected)), strconv.Itoa(game.Game.GamePk)))
}

func InsertPlayByPlayRecords(tx *sql.Tx, gamePk int, records []Play) {
	playByPlayValueStrings := make([]string, 0, len(records))
	playByPlayValueArgs := make([]interface{}, 0, len(records)*17)

	contributorValueStrings := make([]string, 0, len(records))
	contributorValueArgs := make([]interface{}, 0, len(records)*4)

	for _, play := range records {
		playByPlayValueStrings = append(playByPlayValueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		playByPlayValueArgs = append(playByPlayValueArgs,
			gamePk,
			play.About.EventIdx,
			play.About.EvendId,
			play.About.Period,
			play.About.PeriodType,
			play.About.PeriodTime,
			play.About.DateTime,
			play.About.Goals.Away,
			play.About.Goals.Home,
			play.Result.Event,
			play.Result.EventCode,
			play.Result.EventTypeId,
			play.Result.Description,
			play.Result.SecondaryType,
			play.Coordinates.X,
			play.Coordinates.Y,
			play.Team.Id,
		)

		for _, player := range play.Players {
			contributorValueStrings = append(contributorValueStrings, "(?, ?, ?, ?)")
			contributorValueArgs = append(contributorValueArgs,
				gamePk,
				play.About.EventIdx,
				player.Player.PlayerId,
				player.PlayerType,
			)
		}
	}

	stmt := fmt.Sprintf("INSERT INTO play_by_play VALUES %s ON DUPLICATE KEY UPDATE game_pk=game_pk", strings.Join(playByPlayValueStrings, ","))
	result, err := tx.Exec(stmt, playByPlayValueArgs...)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into play_by_play for gamePk %s", strconv.Itoa(int(rows_affected)), strconv.Itoa(gamePk)))

	stmt = fmt.Sprintf("INSERT INTO play_by_play_contributor VALUES %s ON DUPLICATE KEY UPDATE game_pk=game_pk", strings.Join(contributorValueStrings, ","))
	result, err = tx.Exec(stmt, contributorValueArgs...)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	rows_affected, err = result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into play_by_play_contributor for gamePk %s", strconv.Itoa(int(rows_affected)), strconv.Itoa(gamePk)))
}

func InsertOnIceRecords(tx *sql.Tx, records []OnIceRecord) {
	onIceValueStrings := make([]string, 0, len(records))
	onIceValueArgs := make([]interface{}, 0, len(records)*17)

	for _, oir := range records {
		onIceValueStrings = append(onIceValueStrings, "(?, ?, ?, ?, ?)")
		onIceValueArgs = append(onIceValueArgs,
			oir.gamePk,
			oir.teamId,
			oir.eventIdx,
			oir.lineHash,
			oir.goalieId,
		)
	}

	stmt := fmt.Sprintf("INSERT INTO play_by_play_on_ice VALUES %s ON DUPLICATE KEY UPDATE game_pk=game_pk", strings.Join(onIceValueStrings, ","))
	result, err := tx.Exec(stmt, onIceValueArgs...)

	if err != nil {
		tx.Rollback()
		fmt.Println("Rolling back")
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into play_by_play_on_ice for gamePk %s", strconv.Itoa(int(rows_affected)), strconv.Itoa(records[0].gamePk)))
}

func InsertSkaterLineRecords(db *sql.DB, records []OnIceRecord, season int) {
	stakerLineValueStrings := make([]string, 0, len(records))
	stakerLineValueArgs := make([]interface{}, 0, len(records)*17)

	for _, oir := range records {
		stakerLineValueStrings = append(stakerLineValueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		stakerLineValueArgs = append(stakerLineValueArgs,
			season,
			oir.teamId,
			oir.lineHash,
			oir.skaterId1,
			oir.skaterId2,
			oir.skaterId3,
			oir.skaterId4,
			oir.skaterId5,
			oir.skaterId6,
		)
	}

	stmt := fmt.Sprintf("INSERT INTO team_season_skater_lines VALUES %s ON DUPLICATE KEY UPDATE line_hash=line_hash", strings.Join(stakerLineValueStrings, ","))
	result, err := db.Exec(stmt, stakerLineValueArgs...)

	if err != nil {
		fmt.Println("Error inserting staker line record")
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into team_season_skater_lines for gamePk %s", strconv.Itoa(int(rows_affected)), strconv.Itoa(records[0].gamePk)))
}

func UpdateEtlGameStatus(db *sql.DB, gamePk int, status string) {
	stmt := fmt.Sprintf("INSERT INTO etl_game_status VALUES (\"%s\", \"%s\") ON DUPLICATE KEY UPDATE status=\"%s\"", strconv.Itoa(gamePk), status, status)
	result, err := db.Exec(stmt)

	if err != nil {
		fmt.Println("Error updating ETL Game Status")
		log.Fatal(err)
	}

	rows_affected, err := result.RowsAffected()
	println(fmt.Sprintf("Inserted %s records into etl_game_status for gamePk %s", strconv.Itoa(int(rows_affected)), strconv.Itoa(gamePk)))
}

func gameHasBeenProcessed(db *sql.DB, gamePk int) bool {
	query := fmt.Sprintf("SELECT status FROM etl_game_status WHERE game_pk = %s", strconv.Itoa(gamePk))
	results, err := db.Query(query)

	if err != nil {
		log.Fatal("Error retrieving etl game status")
	}

	gameStatus := ""

	for results.Next() {
		err = results.Scan(&gameStatus)
		if err != nil {
			log.Fatal(err)
		}
	}

	return gameStatus == "COMPLETE"
}
