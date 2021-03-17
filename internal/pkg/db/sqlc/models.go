// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
)

type Log struct {
	ID int64 `json:"id"`
	// Anzahl Inserts
	Inserts sql.NullInt32 `json:"inserts"`
	// Anzahl Updates
	Updates sql.NullInt32 `json:"updates"`
	// Anzahl Errors
	Errors sql.NullInt32 `json:"errors"`
	// Startzeitpunkt
	TimestampStarted sql.NullTime `json:"timestamp_started"`
	// Endzeitpunkt
	TimestampFinished sql.NullTime `json:"timestamp_finished"`
}

type Rim struct {
	// Eindeutiger Code
	Code int32 `json:"code"`
	// Felgenmaulweite n Zoll, float mit max 2 vor und 2 nachkomma stellen
	Width sql.NullString `json:"width"`
	// Code für Felgenhöhe
	Height sql.NullString `json:"height"`
	// X = True, x = False
	OnePiece sql.NullBool `json:"one_piece"`
	// Felgendurchmesser in Zoll
	Diameter sql.NullInt32 `json:"diameter"`
	// S = Stahl, L = Leichtmetall
	Material sql.NullString `json:"material"`
}

type Timespan struct {
	// interner Zeitraumschlüssel
	SchwackeID int32 `json:"schwacke_id"`
	// Interner Typschlüssel
	SchwackeCode sql.NullInt32 `json:"schwacke_code"`
	// Datum (TT.MM.JJJJ)
	ValidFrom sql.NullTime `json:"valid_from"`
	// Datum (TT.MM.JJJJ). Wert ist 00.00.0000 falls gültig
	ValidUntil sql.NullTime `json:"valid_until"`
}