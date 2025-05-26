package main

import (
	"database/sql"
	"fmt"

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

// CreateTable creates the evolving_hockey_player_seasons table if it doesn't exist
func (r *PlayerRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS evolving_hockey_player_seasons (
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

// InsertPlayerStats inserts a slice of PlayerStats into the database using a transaction
func (r *PlayerRepository) InsertPlayerStats(players []PlayerStats) error {
	if len(players) == 0 {
		return nil
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare the insert statement
	query := `
	INSERT INTO evolving_hockey_player_seasons 
	(nhl_id, season, full_name, eh_id, team, position, shoots_catches, birthday, 
	 draft_year, draft_round, overall_pick, gp, toi_all, gar, war, spar)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
	updated_at = CURRENT_TIMESTAMP`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert each player record
	for i, player := range players {
		_, err := stmt.Exec(
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
		if err != nil {
			return fmt.Errorf("failed to insert player record %d (%s): %w", i+1, player.FullName, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
