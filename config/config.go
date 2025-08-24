package config

import (
	"backend-assignment/logs"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	logs.Info("Successfully connected to database")
	return db, nil
}
