package models

import (
	"database/sql"
	"time"
)

type Log struct {
	message   string
	timeStamp time.Time
}

type DBLogger struct {
	DB *sql.DB
}

func (logger *DBLogger) Create(err error) {
	//TODO: Implement this so when this function is called we write the log to the Database
}
