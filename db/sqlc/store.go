package db

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mapadj/goSchlacke/util"
)

//Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rberr)
		}
		return err
	}

	return tx.Commit()
}

// ImportTxParams params
type ImportTxParams struct {
	FilePath             string
	Table                string
	DatVersion           string
	MaxFailRateInPerCent int
	MaxErrorCount        int
}

// ImportTxResult params
type ImportTxResult struct {
	NumberOfLines  int
	NumberOfFailes int
	Success        bool
	Inserts        int
	Updates        int
}

// ImportChoserParams params
type ImportChoserParams struct {
	Ctx     context.Context
	Arg     ImportTxParams
	Result  *ImportTxResult
	Scanner *bufio.Scanner
}

// ImportTx import transaction
func (store *Store) ImportTx(ctx context.Context, arg ImportTxParams) (result ImportTxResult, err error) {

	// Init result
	result.NumberOfLines = 0
	result.NumberOfFailes = 0
	result.Success = false
	result.Inserts = 0
	result.Updates = 0

	// Open File
	file, err := os.Open(arg.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Count Lines
	result.NumberOfLines, err = util.LineCounter(file)
	if err != nil {
		return result, err
	}
	println("LineCount: ", result.NumberOfLines)
	// Calculate Max Error Count
	arg.MaxErrorCount = result.NumberOfLines * arg.MaxFailRateInPerCent / 100

	file.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(file)

	importArgs := &ImportChoserParams{
		Ctx:     ctx,
		Arg:     arg,
		Result:  &result,
		Scanner: scanner,
	}

	fn := map[string]func(args *ImportChoserParams) (err error){
		"rimsV1":      store.ConvertRimsV1,
		"timespansV1": store.ConvertTimespansV1,
	}

	// Chose Importfunction and Import
	err = fn[arg.Table+arg.DatVersion](importArgs)

	return result, err
}
