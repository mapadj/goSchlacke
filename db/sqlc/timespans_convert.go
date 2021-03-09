package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

/*
Feldname	 	Typ 	Länge 	Beschreibung
schwacke_id		Num		8		Interner Zeitraumschlüssel
schwacke_code	Num		8		Interner Typschlüssel
valid_from		Date	10		Datum (TT.MM.JJJJ)
valid_to		Date	10		Datum (TT.MM.JJJJ). Wert ist 00.00.0000 falls gültig
*/

// TimespansV1LineLength Length of one Line
const TimespansV1LineLength = 36

// ConvertLineTimespansV1 convert one line
func ConvertLineTimespansV1(line string) (arg UpsertTimespansV1Params, err error) {

	// Check if line size is correct
	if len(line) != TimespansV1LineLength {
		return arg, fmt.Errorf("Line Size wrong")
	}

	// SchwackeId, Num, Len 8
	schwackeIDString := line[0:8]
	schwackeIDInt, err := strconv.Atoi(schwackeIDString)
	if err != nil {
		return arg, fmt.Errorf("SchwackeId: Conversion Error: %v", err)
	}
	arg.SchwackeID = int32(schwackeIDInt)

	// SchwackeId, Num, Len 8
	schwackeCodeString := line[8:16]
	schwackeCodeInt, err := strconv.Atoi(schwackeCodeString)
	if err != nil {
		return arg, fmt.Errorf("SchwackeId: Conversion Error: %v", err)
	}
	arg.SchwackeCode = sql.NullInt32{Int32: int32(schwackeCodeInt), Valid: true}

	// valid_from, Len 10, Datum (TT.MM.JJJJ)
	layout := "02.01.2006"

	validFromString := line[16:26]
	validFromDate, err := time.Parse(layout, validFromString)
	if err != nil {
		return arg, fmt.Errorf("ValidFrom: Conversion Error: %v", err)
	}
	arg.ValidFrom = sql.NullTime{Time: validFromDate, Valid: true}

	// valid_from, Len 10, Datum (TT.MM.JJJJ)
	var validToDate time.Time
	validToString := line[26:36]
	allballs := "00.00.0000"
	if validToString != allballs {
		validToDate, err = time.Parse(layout, validToString)
		if err != nil {
			return arg, fmt.Errorf("ValidFrom: Conversion Error: %v", err)
		}
	}
	arg.ValidFrom = sql.NullTime{Time: validToDate, Valid: true}

	// Return Result
	return arg, nil
}

// ConvertTimespansV1 choses Import
func (store *Store) ConvertTimespansV1(args *ImportChoserParams) (err error) {
	return store.execTx(args.Ctx, func(q *Queries) error {

		// Prepare Scanner
		errorCount := 0
		lineIndex := 0

		// Count Table Size
		numberOfEntries, err := q.CountTimespansV1(args.Ctx)
		if err != nil {
			return err
		}

		// Scan Lines
		for args.Scanner.Scan() {

			// Prepare
			lineIndex++
			line := args.Scanner.Text()

			qarg, err := ConvertLineTimespansV1(line)
			if err != nil {
				errorCount++
				println("line: ", lineIndex, " -> ", err.Error())
				continue
			}
			_, err = q.UpsertTimespansV1(args.Ctx, qarg)
			if err != nil {
				log.Println("Upsert failed", err)
				return err
			}

			if errorCount > args.Arg.MaxErrorCount {
				return errors.New("MaxErrorCount is reached")
			}

		}
		if err := args.Scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Count Table Size
		numberOfEntriesAfter, err := q.CountTimespansV1(args.Ctx)
		if err != nil {
			return err
		}

		// Statistics
		args.Result.NumberOfFailes = errorCount
		args.Result.Inserts = int(numberOfEntriesAfter - numberOfEntries)
		args.Result.Updates = int(args.Result.NumberOfLines - args.Result.NumberOfFailes - args.Result.Inserts)
		println("ErrorCount:", args.Result.NumberOfFailes, "of", args.Result.NumberOfLines)
		args.Result.Success = true

		return nil
	})

}
