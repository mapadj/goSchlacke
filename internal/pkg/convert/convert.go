package convert

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func IsEmpty(i string) bool {
	return len(strings.TrimSpace(i)) == 0
}

func ConvertXxToBool(input string) (bool, error) {
	if input == "X" {
		return true, nil
	} else if input == "x" {
		return false, nil
	} else {
		return false, fmt.Errorf("OnePiece: Wrong Data: %v", input)
	}
}

func ConvertXxToNullBool(input string) (output sql.NullBool, err error) {
	o, err := ConvertXxToBool(input)
	if err != nil {
		return output, err
	}
	return ConvertBoolToNullBool(o), nil
}

func ConvertStringToInt32(input string) (output int32, err error) {
	o, err := strconv.Atoi(input)
	if err != nil {
		return output, fmt.Errorf("Code: Conversion Error: %v", err)
	}
	return int32(o), nil

}

func ConvertStringToNullInt32(input string) (output sql.NullInt32, err error) {
	o, err := ConvertStringToInt32(input)
	if err != nil {
		return output, err
	}
	output = ConvertInt32ToNullInt32(o)
	return output, nil
}

func ConvertStringToNullString(input string) sql.NullString {
	return sql.NullString{String: input, Valid: true}
}
func ConvertInt32ToNullInt32(input int32) sql.NullInt32 {
	return sql.NullInt32{Int32: input, Valid: true}
}
func ConvertBoolToNullBool(input bool) sql.NullBool {
	return sql.NullBool{Bool: input, Valid: true}
}

func ConvertStringToNullDate(input string, format string) (output sql.NullTime, err error) {
	// Parse Value
	validFromDate, err := time.Parse(format, input)
	if err != nil {
		return output, fmt.Errorf("TimeDate: Conversion Error: %v", err)
	}

	// Create NullTime Value
	output = sql.NullTime{Time: validFromDate, Valid: true}

	return output, nil
}

const date_layout_german_10 string = "02.01.2006"

func ConvertStringToNullDateStandardFormatWithAllZeroCheck(input string) (output sql.NullTime, err error) {

	if input == "00.00.0000" {
		return sql.NullTime{Time: time.Time{}, Valid: true}, nil
	}

	return ConvertStringToNullDate(input, date_layout_german_10)
}
