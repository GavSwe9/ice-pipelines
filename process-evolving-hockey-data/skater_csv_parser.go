package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// parseSkaterCSV parses the skaters CSV file and returns a slice of PlayerStats
func parseSkaterCSV(filepath string) ([]PlayerStats, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open skater CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read skater CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("empty skater CSV file")
	}

	// Skip header row
	records = records[1:]

	var players []PlayerStats
	for i, record := range records {
		if len(record) < 25 { // Expected number of columns (reduced by 1 due to removed age)
			return nil, fmt.Errorf("invalid record at line %d: expected 25 columns, got %d", i+2, len(record))
		}

		player, err := parseSkaterRecord(record)
		if err != nil {
			return nil, fmt.Errorf("error parsing record at line %d: %w", i+2, err)
		}

		players = append(players, player)
	}

	return players, nil
}

// parseSkaterRecord converts a CSV record to a PlayerStats struct
func parseSkaterRecord(record []string) (PlayerStats, error) {
	var player PlayerStats
	var err error

	// Remove quotes from string fields
	player.FullName = strings.Trim(record[0], `"`)
	player.EhId = strings.Trim(record[1], `"`)
	player.NhlId = strings.Trim(record[2], `"`)
	player.Season = strings.Trim(record[3], `"`)
	player.Team = strings.Trim(record[4], `"`)
	player.Position = strings.Trim(record[5], `"`)
	player.ShootsCatches = strings.Trim(record[6], `"`)

	// Parse birthday
	birthdayStr := strings.Trim(record[7], `"`)
	player.Birthday, err = time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		return player, fmt.Errorf("failed to parse birthday '%s': %w", birthdayStr, err)
	}

	// Parse draft year (nullable)
	draftYrStr := strings.Trim(record[8], `"`)
	if draftYrStr == "NA" || draftYrStr == "" {
		player.DraftYear = sql.NullInt32{Valid: false}
	} else {
		val, err := strconv.Atoi(draftYrStr)
		if err != nil {
			return player, fmt.Errorf("failed to parse draft year: %w", err)
		}
		player.DraftYear = sql.NullInt32{Int32: int32(val), Valid: true}
	}

	// Parse draft round (nullable)
	draftRdStr := strings.Trim(record[9], `"`)
	if draftRdStr == "NA" || draftRdStr == "" {
		player.DraftRound = sql.NullInt32{Valid: false}
	} else {
		val, err := strconv.Atoi(draftRdStr)
		if err != nil {
			return player, fmt.Errorf("failed to parse draft round: %w", err)
		}
		player.DraftRound = sql.NullInt32{Int32: int32(val), Valid: true}
	}

	// Parse draft overall (nullable)
	draftOvStr := strings.Trim(record[10], `"`)
	if draftOvStr == "NA" || draftOvStr == "" {
		player.OverallPick = sql.NullInt32{Valid: false}
	} else {
		val, err := strconv.Atoi(draftOvStr)
		if err != nil {
			return player, fmt.Errorf("failed to parse draft overall: %w", err)
		}
		player.OverallPick = sql.NullInt32{Int32: int32(val), Valid: true}
	}

	// Parse GP
	player.GP, err = strconv.Atoi(strings.Trim(record[11], `"`))
	if err != nil {
		return player, fmt.Errorf("failed to parse GP: %w", err)
	}

	// Parse TOI_All
	player.ToiAll, err = strconv.ParseFloat(strings.Trim(record[12], `"`), 64)
	if err != nil {
		return player, fmt.Errorf("failed to parse TOI_All: %w", err)
	}

	// Parse GAR (column 21 - index 20)
	player.GAR, err = strconv.ParseFloat(strings.Trim(record[20], `"`), 64)
	if err != nil {
		return player, fmt.Errorf("failed to parse GAR: %w", err)
	}

	// Parse WAR (column 22 - index 21)
	player.WAR, err = strconv.ParseFloat(strings.Trim(record[21], `"`), 64)
	if err != nil {
		return player, fmt.Errorf("failed to parse WAR: %w", err)
	}

	// Parse SPAR (column 23 - index 22)
	player.SPAR, err = strconv.ParseFloat(strings.Trim(record[22], `"`), 64)
	if err != nil {
		return player, fmt.Errorf("failed to parse SPAR: %w", err)
	}

	return player, nil
}
