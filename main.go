package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"time"

	_ "github.com/lib/pq"
	db "github.com/mapadj/goSchlacke/db/sqlc"
	"github.com/mapadj/goSchlacke/util"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/schlacke?sslmode=disable"
)

var queries *db.Queries
var arg *db.ImportTxParams

func init() {
	// Prepare Data
	arg = &db.ImportTxParams{}
	flag.IntVar(&arg.MaxFailRateInPerCent, "max-fail-rate", 5, "Maximal Allowed Failrate in %. Default: 5")
	flag.StringVar(&arg.DatVersion, "version", "V1", "DatVersion, Default: 'V1'")

	flag.StringVar(&arg.FilePath, "file", "", "Path to .dat-File")
	flag.StringVar(&arg.Table, "table", "", "Target Table")

	flag.Parse()

	if arg.FilePath == "" {
		log.Fatalln("file flag missing")
	}

	if arg.Table == "" {
		log.Fatalln("table flag missing")
	}
}

func main() {

	// Track Start Time
	timeStart := time.Now()

	// Load Configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Prepare Database Connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)

	// Import Data
	result, err := store.ImportTx(context.Background(), *arg)
	if err != nil {
		println(err.Error())
	} else {
		println("Inserts: ", result.Inserts, " | Updates: ", result.Updates)
	}

	// Track Finish Time
	timeEnd := time.Now()

	// Prepare Log Data
	argLog := db.InsertLogParams{
		Inserts:           sql.NullInt32{Int32: int32(result.Inserts), Valid: true},
		Updates:           sql.NullInt32{Int32: int32(result.Updates), Valid: true},
		Errors:            sql.NullInt32{Int32: int32(result.NumberOfFailes), Valid: true},
		TimestampStarted:  sql.NullTime{Time: timeStart, Valid: true},
		TimestampFinished: sql.NullTime{Time: timeEnd, Valid: true},
	}

	// Write Log to Database
	_, err = store.InsertLog(context.Background(), argLog)
	if err != nil {
		log.Fatal(err)
	}
}
