package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gavswe19/ice-pipelines/database"
)

func main() {
	db := database.GetDatabase()
	playerIdList := fetchAllPlayers(db, 20212022)

	sqsClient := getSqsClient()
	ctx := context.TODO()

	for _, playerId := range playerIdList {
		fmt.Println(playerId)
		queueMessage := &sqs.SendMessageInput{
			MessageBody: aws.String(strconv.Itoa(playerId)),
			QueueUrl:    aws.String("https://sqs.us-east-1.amazonaws.com/271463937680/ice-player-season-totals-queue"),
		}
		_, err := sqsClient.SendMessage(ctx, queueMessage)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func fetchAllPlayers(db *sql.DB, season int) []int {
	query := fmt.Sprintf("SELECT DISTINCT player_id FROM team_season_players WHERE season = %d", season)
	results, err := db.Query(query)

	if err != nil {
		log.Fatal("Error retrieving team_season_playerss")
	}

	playerIdList := []int{}

	for results.Next() {
		var playerId int
		err = results.Scan(&playerId)
		if err != nil {
			log.Fatal(err)
		}
		playerIdList = append(playerIdList, playerId)
	}
	return playerIdList
}

func getSqsClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)
	return client
}
