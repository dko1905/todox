package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func init() {
	db, err := sql.Open("sqlite3", "sqdata.db")
	if err != nil {
		panic(err)
	}

	database = db
}

func GetConnection(ctx context.Context) (*sql.Conn, error) {
	return database.Conn(ctx)
}
