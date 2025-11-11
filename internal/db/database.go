package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(DSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
