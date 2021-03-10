package db

import (
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
	Ctx            context.Context
	ImportTxParams ImportTxParams
	ImportTxResult *ImportTxResult
	File           *os.File
	Functions      Functions
}

// ImportTx import transaction
func (store *Store) ImportTx(ctx context.Context, importTxParams ImportTxParams) (importTxResult ImportTxResult, err error) {

	// Init result
	importTxResult.NumberOfLines = 0
	importTxResult.NumberOfFailes = 0
	importTxResult.Success = false
	importTxResult.Inserts = 0
	importTxResult.Updates = 0

	// Open File
	file, err := os.Open(importTxParams.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Count Lines
	importTxResult.NumberOfLines, err = util.LineCounter(file)
	if err != nil {
		return importTxResult, err
	}
	// Reset File Pointer
	file.Seek(0, io.SeekStart)

	// Log LineCount
	println("LineCount: ", importTxResult.NumberOfLines)

	// Calculate Max Error Count
	importTxParams.MaxErrorCount = importTxResult.NumberOfLines * importTxParams.MaxFailRateInPerCent / 100

	importChoserParams := &ImportChoserParams{
		Ctx:            ctx,
		ImportTxParams: importTxParams,
		ImportTxResult: &importTxResult,
		File:           file,
	}

	fn := map[string]func(importChoserParams *ImportChoserParams) (err error){
		"rimsV1":      store.ConvertRimsV1,
		"timespansV1": store.ConvertTimespansV1,
	}

	// Chose Importfunction and Import
	err = fn[importTxParams.Table+importTxParams.DatVersion](importChoserParams)

	return importTxResult, err
}

// ConvertRimsV1 choses Import
func (store *Store) ConvertRimsV1(importChoserParams *ImportChoserParams) (err error) {
	return store.execTx(importChoserParams.Ctx, func(q *Queries) error {

		// Count Table Size
		numberOfEntries, err := DBRowCounter[importChoserParams.ImportTxParams.Table+importChoserParams.ImportTxParams.DatVersion](importChoserParams.Ctx)
		if err != nil {
			return err
		}

		errorChannel := make(chan error)

		newLines := FileReader(importChoserParams.File)
		primitveStructs := StructPress(newLines)
		convertedStructs := StartConvertAndValidateThread(primitveStructs, *importChoserParams, errorChannel)

	enough:
		// This is the Query Loop. It waits for receives Data from channel X until it closes.
		for {
			select {
			case err := <-errorChannel:
				return err
			// Wait for data
			case val, ok := <-convertedStructs:
				// check data health
				if !ok {
					log.Println(val, ok, "loop broke")
					break enough
				}

				// Add query to transaction
				_, err = q.UpsertRimsV1(importChoserParams.Ctx, val)
				if err != nil {
					log.Println("Upsert failed: ", err)
					return err
				}
			}
		}

		// Count Table Size
		numberOfEntriesAfter, err := q.CountRimsV1(importChoserParams.Ctx)
		if err != nil {
			return err
		}

		// Statistics
		r := importChoserParams.ImportTxResult
		r.Inserts = int(numberOfEntriesAfter - numberOfEntries)
		r.Updates = int(r.NumberOfLines - r.NumberOfFailes - r.Inserts)
		r.Success = true
		return nil
	})

}
