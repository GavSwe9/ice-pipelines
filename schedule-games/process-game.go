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

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

type GameResponse struct {
	LiveData LiveData `json:"liveData"`
}

type LiveData struct {
	Plays Plays `json:"plays"`
}

type Plays struct {
	AllPlays []Play `json:"allPlays"`
}

type Play struct {
	Players []Players `json:"players"`
	Result Result `json:"result"`
	About About `json:"about"`
	Coordinates Coordinates `json:"coordinates"`
	Team Team `json:"team"`
}

type Players struct {
	Player Player `json:"player"`
	PlayerType string `json:"playerType"`
}

type Player struct {
	PlayerId int `json:"id"`
}

type Result struct {
	Event string `json:"event"`
	EventCode string `json:"eventCode"`
	EventTypeId string `json:"eventTypeId"`
	Description string `json:"description"`
	SecondaryType string `json:"string"`
}

type About struct {
	EventIdx int `json:"eventIdx"`
	EvendId int `json:"eventId"`
	Period int `json:"period"`
	PeriodType string `json:"periodType"`
	OrdinalNum string `json:"ordinalNum"`
	PeriodTime string `json:"periodTime"`
	PeriodTimeRemaining string `json:"periodTimeRemaining"`
	DateTime string `json:"dateTime"`
	Goals Goals `json:"goals"`
}

type Coordinates struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Team struct {
	Id float32 `json:"id"`
}

type Goals struct {
	Away int `json:"away"`
	Home int `json:"home"`
}

var db *sql.DB

// Handler is our lambda handler invoked by the `lambda.Start` function call
func ProcessGame(gamePk int) {
	fmt.Println(fmt.Sprintf(" *** Processing GamePk %s ***", strconv.Itoa(gamePk)));
	response, err := http.Get(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/game/%s/feed/live", strconv.Itoa(gamePk)));

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	
	if err != nil {
		log.Fatal(err)
	}
	
	var responseObject GameResponse
	json.Unmarshal(responseData, &responseObject)
	
	secrets := getAwsSecrets()
	
	cfg := mysql.Config{
		User:   secrets.Username,
		Passwd: secrets.Password,
		Net:    "tcp",
		Addr:   "farm.cxqsjcdo8n1w.us-east-1.rds.amazonaws.com",
		DBName: "ICE",
		AllowNativePasswords: true,
	}
	
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	
	tx, err := db.Begin() 
	println("Start Transation")

	DeleteWithGamePk(tx, "play_by_play", gamePk);
	DeleteWithGamePk(tx, "play_by_play_contributor", gamePk);
	InsertRecords(tx, gamePk, responseObject.LiveData.Plays.AllPlays);

	err = tx.Commit()
	if err != nil {
		log.Fatal(err);
	}
	println("Committed Transation")
}

func DeleteWithGamePk(tx *sql.Tx, tableName string, gamePk int) {
	result, err := tx.Exec(fmt.Sprintf("DELETE FROM %s PBP WHERE PBP.game_pk = %s", tableName, strconv.Itoa(gamePk)));
	if err != nil {
		tx.Rollback();
		log.Fatal(err);
	}

	rows_affected, err := result.RowsAffected();
	println(fmt.Sprintf("Delted %s records from %s for gamePk %s", strconv.Itoa(int(rows_affected)), tableName, strconv.Itoa(gamePk)));	
}

func InsertRecords(tx *sql.Tx, gamePk int, records []Play) {
	playByPlayValueStrings := make([]string, 0, len(records))
	playByPlayValueArgs := make([]interface{}, 0, len(records) * 17)

	contributorValueStrings := make([]string, 0, len(records))
	contributorValueArgs := make([]interface{}, 0, len(records) * 4)

	for _, play := range records {
		playByPlayValueStrings = append(playByPlayValueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
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
			);
		
		for _, player := range play.Players {
			contributorValueStrings = append(contributorValueStrings, "(?, ?, ?, ?)");
			contributorValueArgs = append(contributorValueArgs, 
				gamePk,
				play.About.EventIdx, 
				player.Player.PlayerId,
				player.PlayerType,
				) 
		}
	}

	stmt := fmt.Sprintf("INSERT INTO play_by_play VALUES %s", strings.Join(playByPlayValueStrings, ","))
	result, err := tx.Exec(stmt, playByPlayValueArgs...)

	if err != nil {
		tx.Rollback();
		log.Fatal(err);
	}

	rows_affected, err := result.RowsAffected();
	println(fmt.Sprintf("Inserted %s records into play_by_play for gamePk %s", strconv.Itoa(int(rows_affected)), strconv.Itoa(gamePk)));

	stmt = fmt.Sprintf("INSERT INTO play_by_play_contributor VALUES %s", strings.Join(contributorValueStrings, ","))
	result, err = tx.Exec(stmt, contributorValueArgs...)

	if err != nil {
		tx.Rollback();
		log.Fatal(err);
	}

	rows_affected, err = result.RowsAffected();
	println(fmt.Sprintf("Inserted %s records into play_by_play_contributor for gamePk %s", strconv.Itoa(int(rows_affected)), strconv.Itoa(gamePk)));	
}
