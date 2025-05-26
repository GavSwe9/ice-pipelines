package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gavswe19/ice-pipelines/database"
)

// PlayerRepository handles database operations for player statistics
type PlayerRepository struct {
	db *sql.DB
}

// NewPlayerRepository creates a new PlayerRepository with a database connection
func NewPlayerRepository() (*PlayerRepository, error) {
	db := database.GetDatabase()

	return &PlayerRepository{db: db}, nil
}

// CreateTable creates the evolving_hockey_player_seasons_gar table if it doesn't exist
func (r *PlayerRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS evolving_hockey_player_seasons_gar (
		nhl_id VARCHAR(255) NOT NULL,
		season VARCHAR(10) NOT NULL,
		full_name VARCHAR(255) NOT NULL,
		eh_id VARCHAR(255) NOT NULL,
		team VARCHAR(10) NOT NULL,
		position VARCHAR(10) NOT NULL,
		shoots_catches VARCHAR(5) NOT NULL,
		birthday DATE NOT NULL,
		draft_year INT NULL,
		draft_round INT NULL,
		overall_pick INT NULL,
		gp INT NOT NULL,
		toi_all DECIMAL(10,2) NOT NULL,
		gar DECIMAL(10,2) NOT NULL,
		war DECIMAL(10,2) NOT NULL,
		spar DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (nhl_id, season),
		INDEX idx_team_season (team, season)
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// InsertPlayerStats inserts a slice of PlayerStats into the database using batch inserts
func (r *PlayerRepository) InsertPlayerStats(players []PlayerStats) error {
	if len(players) == 0 {
		return nil
	}

	const batchSize = 1000 // Insert 1000 records at a time

	// Process in batches
	for i := 0; i < len(players); i += batchSize {
		end := i + batchSize
		if end > len(players) {
			end = len(players)
		}

		batch := players[i:end]
		if err := r.insertBatch(batch); err != nil {
			return fmt.Errorf("failed to insert batch starting at index %d: %w", i, err)
		}
	}

	return nil
}

// insertBatch inserts a batch of players using a single query
func (r *PlayerRepository) insertBatch(players []PlayerStats) error {
	if len(players) == 0 {
		return nil
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Build the batch insert query
	valueStrings := make([]string, 0, len(players))
	valueArgs := make([]interface{}, 0, len(players)*16) // 16 columns per player

	for _, player := range players {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			player.NhlId,
			player.Season,
			player.FullName,
			player.EhId,
			player.Team,
			player.Position,
			player.ShootsCatches,
			player.Birthday,
			player.DraftYear,
			player.DraftRound,
			player.OverallPick,
			player.GP,
			player.ToiAll,
			player.GAR,
			player.WAR,
			player.SPAR,
		)
	}

	query := fmt.Sprintf(`
	INSERT INTO evolving_hockey_player_seasons_gar 
	(nhl_id, season, full_name, eh_id, team, position, shoots_catches, birthday, 
	 draft_year, draft_round, overall_pick, gp, toi_all, gar, war, spar)
	VALUES %s
	ON DUPLICATE KEY UPDATE
	full_name = VALUES(full_name),
	eh_id = VALUES(eh_id),
	team = VALUES(team),
	position = VALUES(position),
	shoots_catches = VALUES(shoots_catches),
	birthday = VALUES(birthday),
	draft_year = VALUES(draft_year),
	draft_round = VALUES(draft_round),
	overall_pick = VALUES(overall_pick),
	gp = VALUES(gp),
	toi_all = VALUES(toi_all),
	gar = VALUES(gar),
	war = VALUES(war),
	spar = VALUES(spar),
	updated_at = CURRENT_TIMESTAMP`, strings.Join(valueStrings, ","))

	// Execute the batch insert
	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
