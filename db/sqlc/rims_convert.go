package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*
Feldname 	Typ 	Länge 	Beschreibung

code 		Num 	5 		Eindeutiger Code
width 		Num		5		Felgenmaulweite n Zoll
height		Char	1		Code für Felgenhöhe
one_piece	Bool	1		X = True, x = False
diameter	Num		2		Felgendurchmesser in Zoll
material	Char	1		S = Stahl, L = Leichtmetall
*/

// RimsV1LineLength Length of one Line
const RimsV1LineLength = 15

// ConvertLineRimsV1 convert one line
func ConvertLineRimsV1(line string) (arg UpsertRimsV1Params, err error) {

	// Check if line size is correct
	if len(line) != RimsV1LineLength {
		return arg, fmt.Errorf("Line Size wrong")
	}

	// Code
	codeString := line[0:5]
	codeInt, err := strconv.Atoi(codeString)
	if err != nil {
		return arg, fmt.Errorf("Code: Conversion Error: %v", err)
	}
	arg.Code = int32(codeInt)

	// Width
	widthString := strings.TrimSpace(line[5:10])
	// Check if String can be converted to Float
	_, err = strconv.ParseFloat(widthString, 64)
	if err != nil {
		return arg, fmt.Errorf("Width: Conversion Error: %v", err)
	}
	arg.Width = sql.NullString{String: widthString, Valid: true}

	// Height Char
	heightString := line[10:11]
	possibleHeights := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if heightString == " " {
		return arg, fmt.Errorf("Height: Missing Data")
	}
	if ok := strings.Contains(possibleHeights, heightString); !ok {
		return arg, fmt.Errorf("Height: Wrong Data")
	}
	arg.Height = sql.NullString{String: heightString, Valid: true}

	// One Piece Bool
	onePieceString := line[11:12]
	var onePieceBool bool
	if onePieceString == "X" {
		onePieceBool = true
	} else if onePieceString == "x" {
		onePieceBool = false
	} else {
		return arg, fmt.Errorf("OnePiece: Wrong Data: %v", onePieceString)
	}
	arg.OnePiece = sql.NullBool{Bool: onePieceBool, Valid: true}

	// Diameter
	diameterString := line[12:14]
	diameterInt, err := strconv.Atoi(diameterString)
	if err != nil {
		return arg, fmt.Errorf("Diameter: Wrong Data %v", err)
	}
	arg.Diameter = sql.NullInt32{Int32: int32(diameterInt), Valid: true}

	// Material Character
	materialString := line[14:15]
	materials := "SL"
	if !strings.Contains(materials, materialString) {
		return arg, fmt.Errorf("Material: WrongData: %v", materialString)
	}
	arg.Material = sql.NullString{String: materialString, Valid: true}

	return arg, nil
}

// ConvertRimsV1 choses Import
func (store *Store) ConvertRimsV1(args *ImportChoserParams) (err error) {
	return store.execTx(args.Ctx, func(q *Queries) error {

		// Prepare Scanner
		errorCount := 0
		lineIndex := 0

		// Count Table Size
		numberOfEntries, err := q.CountRimsV1(args.Ctx)
		if err != nil {
			return err
		}

		// Scan Lines
		for args.Scanner.Scan() {

			// Prepare
			lineIndex++
			line := args.Scanner.Text()

			qarg, err := ConvertLineRimsV1(line)
			if err != nil {
				errorCount++
				println("line: ", lineIndex, " -> ", err.Error())
				continue
			}
			_, err = q.UpsertRimsV1(args.Ctx, qarg)
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
		numberOfEntriesAfter, err := q.CountRimsV1(args.Ctx)
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
