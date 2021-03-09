package db

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomRimCode generates a random Schwacke Rim Code
func RandomRimCode() int32 {
	return RandomInt(0, 99999)
}

const rimMaterials = "SL"

// RandomRimMaterial generates a random rim Material
func RandomRimMaterial() sql.NullString {
	// Create String Builder
	var sb strings.Builder

	// generate Byte
	c := rimMaterials[rand.Intn(len(rimMaterials))]

	// Write Byte to String Builder
	sb.WriteByte(c)

	// Generate String
	randomString := sb.String()

	return sql.NullString{String: randomString, Valid: true}
}

// RandomRimV1 generate Random Schlacke RimV1 Data
func RandomRimV1() UpsertRimsV1Params {
	return UpsertRimsV1Params{
		Code:     RandomRimCode(), // Example for Random Code Generation
		Width:    sql.NullString{String: "12.50", Valid: true},
		Height:   sql.NullString{String: "A", Valid: true},
		OnePiece: sql.NullBool{Bool: true, Valid: true},
		Diameter: sql.NullInt32{Int32: 12, Valid: true},
		Material: RandomRimMaterial(), // Example for Random Material Generation
	}
}
