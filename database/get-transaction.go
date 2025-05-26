package database

import (
	"database/sql"
	"log"

	aws_secrets "github.com/gavswe19/ice-pipelines/aws-secrets"
	"github.com/go-sql-driver/mysql"
)

func GetTransaction() *sql.Tx {
	secrets := aws_secrets.GetAwsSecrets()

	cfg := mysql.Config{
		User:                 secrets.Username,
		Passwd:               secrets.Password,
		Net:                  "tcp",
		Addr:                 "farm.cxqsjcdo8n1w.us-east-1.rds.amazonaws.com",
		DBName:               "ICE",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	return tx
}
