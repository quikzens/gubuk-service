package db

import (
	"database/sql"
	"log"

	"gubuk-service/config"
	sqlc "gubuk-service/db/sqlc"

	_ "github.com/lib/pq"
)

var (
	DB      *sql.DB
	Queries *sqlc.Queries
)

func init() {
	db, err := sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	DB = db
	Queries = sqlc.New(db)
}
