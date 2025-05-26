package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gavswe19/ice-pipelines/database"
)

// Define the Go struct to match the JSON structure
type Player struct {
	PlayerId            int         `json:"playerId"`
	IsActive            bool        `json:"isActive"`
	CurrentTeamId       int         `json:"currentTeamId"`
	CurrentTeamAbbrev   string      `json:"currentTeamAbbrev"`
	FullTeamName        Translation `json:"fullTeamName"`
	FirstName           Translation `json:"firstName"`
	LastName            Translation `json:"lastName"`
	TeamLogo            string      `json:"teamLogo"`
	SweaterNumber       int         `json:"sweaterNumber"`
	Position            string      `json:"position"`
	Headshot            string      `json:"headshot"`
	HeroImage           string      `json:"heroImage"`
	HeightInInches      int         `json:"heightInInches"`
	HeightInCentimeters int         `json:"heightInCentimeters"`
	WeightInPounds      int         `json:"weightInPounds"`
	WeightInKilograms   int         `json:"weightInKilograms"`
	BirthDate           string      `json:"birthDate"`
	BirthCity           Translation `json:"birthCity"`
	BirthStateProvince  Translation `json:"birthStateProvince"`
	BirthCountry        string      `json:"birthCountry"`
	ShootsCatches       string      `json:"shootsCatches"`
	DraftDetails        Draft       `json:"draftDetails"`
}

// Define a struct for fields that have translations
type Translation struct {
	Default string `json:"default"`
	Fr      string `json:"fr,omitempty"` // Use 'omitempty' for optional fields
}

// Define the Draft struct for the draftDetails field
type Draft struct {
	Year        int    `json:"year"`
	TeamAbbrev  string `json:"teamAbbrev"`
	Round       int    `json:"round"`
	PickInRound int    `json:"pickInRound"`
	OverallPick int    `json:"overallPick"`
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	eventRecord := sqsEvent.Records[0]
	eventBody := eventRecord.Body
	playerId, err := strconv.Atoi(eventBody)

	if err != nil {
		fmt.Println("Error handling SQS message:", err)
		return
	}

	// Replace this with the actual API URL you want to use
	apiURL := fmt.Sprintf("https://api-web.nhle.com/v1/player/%d/landing", playerId)

	// Make the HTTP GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON data into the Go struct
	var player Player
	err = json.Unmarshal(body, &player)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Print the struct to verify the data
	fmt.Printf("%+v\n", player)

	insertPlayer(player)
}

func insertPlayer(player Player) {
	query := `
	INSERT INTO player_bio (
		player_id, is_active, current_team_id, current_team_abbrev, full_team_name, first_name, last_name, full_name, sweater_number,
		position, headshot, hero_image, height_in_inches, height_in_centimeters, weight_in_pounds, weight_in_kilograms, birth_date,
		birth_city, birth_state_province, birth_country, shoots_catches, draft_year, draft_team_abbrev, draft_round, draft_pick_in_round,
		draft_overall_pick
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
		is_active = VALUES(is_active),
		current_team_id = VALUES(current_team_id),
		current_team_abbrev = VALUES(current_team_abbrev),
		full_team_name = VALUES(full_team_name),
		first_name = VALUES(first_name),
		last_name = VALUES(last_name),
		full_name = VALUES(full_name),
		sweater_number = VALUES(sweater_number),
		position = VALUES(position),
		headshot = VALUES(headshot),
		hero_image = VALUES(hero_image),
		height_in_inches = VALUES(height_in_inches),
		height_in_centimeters = VALUES(height_in_centimeters),
		weight_in_pounds = VALUES(weight_in_pounds),
		weight_in_kilograms = VALUES(weight_in_kilograms),
		birth_date = VALUES(birth_date),
		birth_city = VALUES(birth_city),
		birth_state_province = VALUES(birth_state_province),
		birth_country = VALUES(birth_country),
		shoots_catches = VALUES(shoots_catches),
		draft_year = VALUES(draft_year),
		draft_team_abbrev = VALUES(draft_team_abbrev),
		draft_round = VALUES(draft_round),
		draft_pick_in_round = VALUES(draft_pick_in_round),
		draft_overall_pick = VALUES(draft_overall_pick)
	 `

	db := database.GetDatabase()
	_, err := db.Exec(
		query,
		player.PlayerId,
		player.IsActive,
		player.CurrentTeamId,
		player.CurrentTeamAbbrev,
		player.FullTeamName.Default, // Assuming using the default language for FullTeamName
		player.FirstName.Default,    // Assuming using the default language for FirstName
		player.LastName.Default,     // Assuming using the default language for LastName
		fmt.Sprintf("%s %s", player.FirstName.Default, player.LastName.Default),
		player.SweaterNumber,
		player.Position,
		player.Headshot,
		player.HeroImage,
		player.HeightInInches,
		player.HeightInCentimeters,
		player.WeightInPounds,
		player.WeightInKilograms,
		player.BirthDate,
		player.BirthCity.Default,          // Assuming using the default language for BirthCity
		player.BirthStateProvince.Default, // Assuming using the default language for BirthStateProvince
		player.BirthCountry,
		player.ShootsCatches,
		player.DraftDetails.Year,
		player.DraftDetails.TeamAbbrev,
		player.DraftDetails.Round,
		player.DraftDetails.PickInRound,
		player.DraftDetails.OverallPick,
	)

	if err != nil {
		log.Fatalf("Failed to insert player: %v", err)
	} else {
		fmt.Println("Player inserted successfully!")
	}
}
