package database

import (
	"database/sql"
	"log"

	aws_secrets "github.com/gavswe19/ice-pipelines/aws-secrets"
	"github.com/go-sql-driver/mysql"
)

func GetDatabase() *sql.DB {
	secrets := aws_secrets.GetAwsSecrets()

	cfg := mysql.Config{
		User:                 secrets.Username,
		Passwd:               secrets.Password,
		Net:                  "tcp",
		Addr:                 "farm.cxqsjcdo8n1w.us-east-1.rds.amazonaws.com",
		DBName:               "ICE",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db
}
