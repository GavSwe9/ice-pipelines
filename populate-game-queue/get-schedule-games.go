package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

type ScheduleResponse struct {
	Dates []Dates `json:"dates"`
}

type Dates struct {
	Games []Game `json:"games"`
}

type Game struct {
	GamePk   int    `json:"gamePk"`
	GameType string `json:"gameType"`
	Season   string `json:"season"`
	GameDate string `json:"gameDate"`
	Teams    Teams  `json:"teams"`
}

type Teams struct {
	AwayTeam GameTeam `json:"away"`
	HomeTeam GameTeam `json:"home"`
}

type GameTeam struct {
	Team Team `json:"team"`
}

type Team struct {
	Id int `json:"id"`
}

// Returns all GamePks for the given date
func GetScheduleGames(dte time.Time) (gamePkList []int) {
	// currentTime := time.Now()
	// dateStr := currentTime.Format("2006-01-02")
	// response, err := http.Get(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/schedule?startDate=%s&endDate=%s", dateStr, dateStr));
	dateStr := dte.Format("2006-01-02")
	response, err := http.Get(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/schedule?startDate=%s&endDate=%s", dateStr, dateStr))

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responseObject ScheduleResponse
	json.Unmarshal(responseData, &responseObject)

	if len(responseObject.Dates) == 0 {
		return
	}

	gamePkList = make([]int, 0, len(responseObject.Dates[0].Games))
	for _, game := range responseObject.Dates[0].Games {
		if game.GameType != "PR" {
			gamePkList = append(gamePkList, game.GamePk)
		}
	}
	return
}
