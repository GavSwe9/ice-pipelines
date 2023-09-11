package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
)

type GatewayResponse events.APIGatewayProxyResponse

func main() {
	// PopulateQueue()
	runYear()
}

func runYear() {
	dte := time.Date(2021, 10, 1, 0, 0, 0, 0, time.Local)
	stopDate := time.Date(2022, 8, 1, 0, 0, 0, 0, time.Local)

	for dte.UnixMilli() < stopDate.UnixMilli() {
		fmt.Println(dte.Format("2006-01-02"))
		PopulateQueue(dte)
		dte = dte.AddDate(0, 0, 1)
	}
}

func PopulateQueue(dte time.Time) {
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
