package db

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBData struct {
	DBCon  string
	DBHost string
	DBName string
	DBUser string
	DBPass string
	DBPort string
}

func GetDBClient() *sqlx.DB {
	var dbCreds = DBData{}
	setDBData(&dbCreds)

	db, err := sqlx.Open(
		dbCreds.DBCon,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbCreds.DBUser,
			dbCreds.DBPass,
			dbCreds.DBHost,
			dbCreds.DBPort,
			dbCreds.DBName,
		),
	)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func setDBData(db *DBData) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "banking"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPass := os.Getenv("DB_PASSWORD")

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3307"
	}

	dbCon := os.Getenv("DB_CON")
	if dbCon == "" {
		dbCon = "mysql"
	}

	db.DBHost = dbHost
	db.DBName = dbName
	db.DBUser = dbUser
	db.DBPass = dbPass
	db.DBPort = dbPort
	db.DBCon = dbCon
}
