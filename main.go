package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"time"

	_ "github.com/lib/pq"
	db "github.com/mapadj/goSchlacke/db/sqlc"
	"github.com/mapadj/goSchlacke/db/sqlc/rims"
	"github.com/mapadj/goSchlacke/util"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/schlacke?sslmode=disable"
)

var queries *db.Queries
var arg *db.ImportTxParams

type Functions struct {
	ImportTableFactory   func() interface{}
	CountDatabaseEntries func(ctx context.Context) (count int64, err error)
	ConvertAndValidate   func(importTable interface{}) (result interface{}, err error)
	Upsert               func(ctx context.Context, params interface{}) (result interface{}, err error)
}

var FunctionsRimsV1 = Functions{
	ImportTableFactory:   rims.NewImportTable(),
	CountDatabaseEntries: queries.CountRimsV1,
	ConvertAndValidate:   rims.ConvertAndValidate,
	Upsert:               db.UpsertRimsV1Params,
}

var DBRowCounter = &map[string]func(ctx context.Context) (count int64, err error){
	"RimsV1":     queries.CountRimsV1,
	"TimespanV1": queries.CountTimespansV1,
}

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
		println("Importing Data Failed: ", err.Error())
	} else {
		println("Inserts: ", result.Inserts, " | Updates: ", result.Updates, " | Errors: ", result.NumberOfFailes, " | %: ", result.NumberOfFailes/result.NumberOfLines)
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
