package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type GatewayResponse events.APIGatewayProxyResponse

// func Handler(ctx context.Context) {
func Handler() {
	gamePkList := GetScheduleGames();
	fmt.Print(gamePkList)

	for _, gamePk := range gamePkList {
		ProcessGame(gamePk);
	}
}

func main() {
	// Handler()
	// lambda.Start(Handler)
	ProcessGame(2022020289)
}
