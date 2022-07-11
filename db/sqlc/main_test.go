package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sRRRs-7/GachaPon/utils"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Load config error: ", err)
	}

	testDB, err = sql.Open(config.DBdriver, config.DBsource)
	if err != nil {
		log.Fatal("Open database error: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}