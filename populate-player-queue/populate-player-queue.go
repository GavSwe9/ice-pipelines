package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	// "log"
	// "strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	teamAbvList := getAllTeamAbv()
	fmt.Println(teamAbvList)
	fmt.Println(len(teamAbvList))

	sqsClient := getSqsClient()
	ctx := context.TODO()

	for _, teamAbv := range teamAbvList {
		fmt.Println(teamAbv, " ----------------------- ")
		playerIdList := getTeamPlayerIdList(teamAbv)
		for _, playerId := range playerIdList {
			fmt.Println(playerId)
			queueMessage := &sqs.SendMessageInput{
				MessageBody: aws.String(strconv.Itoa(playerId)),
				QueueUrl:    aws.String("https://sqs.us-east-1.amazonaws.com/271463937680/player-queue"),
			}
			_, err := sqsClient.SendMessage(ctx, queueMessage)

			if err != nil {
				log.Fatal(err)
			}
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
