package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// conn, err := sql.Open(dbDriver, dbSource)
	testDB1, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Failed to open database", err)
	}

	testQueries = New(testDB1)
	testDB = testDB1

	os.Exit(m.Run())
}
