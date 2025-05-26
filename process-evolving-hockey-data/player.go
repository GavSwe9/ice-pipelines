package main

import (
	"database/sql"
	"time"
)

// PlayerStats represents a player's statistics for a single season
type PlayerStats struct {
	NhlId         string        `db:"nhl_id"`
	Season        string        `db:"season"`
	FullName      string        `db:"full_name"`
	EhId          string        `db:"eh_id"`
	Team          string        `db:"team"`
	Position      string        `db:"position"`
	ShootsCatches string        `db:"shoots_catches"`
	Birthday      time.Time     `db:"birthday"`
	DraftYear     sql.NullInt32 `db:"draft_year"`
	DraftRound    sql.NullInt32 `db:"draft_round"`
	OverallPick   sql.NullInt32 `db:"overall_pick"`
	GP            int           `db:"gp"`
	ToiAll        float64       `db:"toi_all"`
	GAR           float64       `db:"gar"`
	WAR           float64       `db:"war"`
	SPAR          float64       `db:"spar"`
}
