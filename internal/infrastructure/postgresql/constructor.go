package postgresql

import (
	"database/sql"
	"fmt"
	"os"
)



type Psql struct {
	DB *sql.DB
}

var (
	PSQL_HOST     = os.Getenv("PSQL_HOST")
	PSQL_PORT     = os.Getenv("PSQL_PORT")
	PSQL_USER     = os.Getenv("PSQL_USER")
	PSQL_PASS     = os.Getenv("PSQL_PASS")
	PSQL_DATABASE = os.Getenv("PSQL_DATABASE")
)



var instance *sql.DB

func NewInstance() *sql.DB {

	if instance == nil {
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			PSQL_HOST, PSQL_PORT, PSQL_USER, PSQL_PASS, PSQL_DATABASE)

		db, err := sql.Open("postgres", psqlInfo)

		if err != nil {
			panic(db)
		}

		instance = db
	}

	return instance
}

