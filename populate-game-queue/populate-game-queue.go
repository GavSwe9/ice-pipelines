package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type GatewayResponse events.APIGatewayProxyResponse

func main() {
	yesterday := time.Now().AddDate(0, 0, -1)
	PopulateGameQueue(yesterday)
	// runYear()
}

func runYear() {
	dte := time.Date(2021, 10, 1, 0, 0, 0, 0, time.Local)
	stopDate := time.Date(2022, 8, 1, 0, 0, 0, 0, time.Local)

	for dte.UnixMilli() < stopDate.UnixMilli() {
		fmt.Println(dte.Format("2006-01-02"))
		PopulateGameQueue(dte)
		dte = dte.AddDate(0, 0, 1)
	}
}

func PopulateGameQueue(dte time.Time) {
	gamePkList := GetScheduleGames(dte)
	fmt.Println(gamePkList)

	sqsClient := getSqsClient()
	ctx := context.TODO()

	for _, gamePk := range gamePkList {
		queueMessage := &sqs.SendMessageInput{
			MessageBody: aws.String(strconv.Itoa(gamePk)),
			QueueUrl:    aws.String("https://sqs.us-east-1.amazonaws.com/271463937680/ice-game-queue"),
		}
		_, err := sqsClient.SendMessage(ctx, queueMessage)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func getSqsClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)
	return client
}
