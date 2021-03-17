package SchlackeImporter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	db "github.com/mapadj/goSchlacke/internal/pkg/db/sqlc"
	"github.com/mapadj/goSchlacke/internal/pkg/tables"
	"github.com/mapadj/goSchlacke/internal/pkg/util"
	"github.com/mapadj/goSchlacke/internal/pkg/worker"
)

//Store provides all functions to execute db queries and transactions
type Store struct {
	*db.Queries
	db *sql.DB
}

// NewStore creates new Store
func NewStore(database *sql.DB) *Store {
	return &Store{
		db:      database,
		Queries: db.New(database),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rberr)
		}
		return err
	}

	return tx.Commit()
}

// ImportTx import transaction
func (store *Store) ImportTx(ctx context.Context, importTxParams tables.ImportTxParams) (importTxResult tables.ImportTxResult, err error) {

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

	importChoserParams := &tables.ImportChoserParams{
		Ctx:            ctx,
		ImportTxParams: importTxParams,
		ImportTxResult: &importTxResult,
		File:           file,
	}

	// Chose Importfunction and Import
	err = store.StartImportTransaction(importChoserParams)

	return importTxResult, err
}

// StartImportTransaction choses Import
func (store *Store) StartImportTransaction(importChoserParams *tables.ImportChoserParams) (err error) {
	return store.execTx(importChoserParams.Ctx, func(q *db.Queries) (err error) {

		log.Println("Transaction Started")
		// construct Type Name
		t := strings.Title(importChoserParams.ImportTxParams.Table)
		v := strings.Title(importChoserParams.ImportTxParams.DatVersion)
		n := t + v

		log.Println("n: ", n)
		// select ImportableContainerFactory
		factory := ContainerFactoryMap[n]()

		// Generic Database Count Call
		countMeth := reflect.ValueOf(q).MethodByName("Count" + n)
		resultSlice := countMeth.Call([]reflect.Value{reflect.ValueOf(importChoserParams.Ctx)})
		if len(resultSlice) != 2 {
			return errors.New("reflect count database call returned wrong result length")
		}
		numberOfEntries := resultSlice[0].Interface().(int64)
		rerr := resultSlice[1].Interface()
		if rerr != nil {
			return fmt.Errorf("error calling: %v", rerr)
		}

		// Create Error Channel
		errorChannel := make(chan error)

		// Start File Reader
		newLines := worker.StartFileReader(importChoserParams.File)

		// Convert Byte Lines into Struct using fixed size tables
		primitiveStructs := worker.StructPress(newLines, factory)

		// Convert and Validate Data into Upsertable Structs
		convertedStructs := worker.StartConvertAndValidateThread(primitiveStructs, *importChoserParams, errorChannel)

		// Get Upset Method
		upsertMethod := reflect.ValueOf(q).MethodByName("Upsert" + n)

		err = worker.Upserter(errorChannel, convertedStructs, upsertMethod, importChoserParams.Ctx)
		if err != nil {
			return err
		}

		// Generic Database Count Call Again after Upserts
		resultSlice = countMeth.Call([]reflect.Value{reflect.ValueOf(importChoserParams.Ctx)})
		if len(resultSlice) != 2 {
			return errors.New("reflect count database call returned wrong result length")
		}
		numberOfEntriesAfter := resultSlice[0].Interface().(int64)
		rerr = resultSlice[1].Interface()
		if rerr != nil {
			return fmt.Errorf("error calling: %v", rerr)
		}

		// Statistics
		r := importChoserParams.ImportTxResult
		r.Inserts = int(numberOfEntriesAfter - numberOfEntries)
		r.Updates = int(r.NumberOfLines - r.NumberOfFailes - r.Inserts)
		r.Success = true
		return nil
	})

}
