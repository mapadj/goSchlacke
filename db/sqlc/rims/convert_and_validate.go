package rims

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	db "github.com/mapadj/goSchlacke/db/sqlc"
)

// ConvertAndValidate Convert Function
func ConvertAndValidate(rim ImportTable) (r db.UpsertRimsV1Params, err error) {
	// CODE
	if db.IsEmpty(rim.Code) {
		return r, errors.New("Code Field Empty")
	}
	r.Code, err = db.ConvertStringToInt32(rim.Code)
	if err != nil {
		return r, err
	}

	// WIDTH
	if db.IsEmpty(rim.Width) {
		return r, errors.New("Width Field Empty")
	}
	w := strings.TrimSpace(rim.Width)
	_, err = strconv.ParseFloat(w, 64)
	if err != nil {
		return r, fmt.Errorf("Width: Float Conversion Error: %v", err)
	}
	r.Width = db.ConvertStringToNullString(w)
	if err != nil {
		return r, err
	}

	// OnePiece
	if db.IsEmpty(rim.OnePiece) {
		return r, errors.New("OnePiece Field Empty")
	}
	r.OnePiece, err = db.ConvertXxToNullBool(rim.OnePiece)
	if err != nil {
		return r, err
	}

	// Diameter
	if db.IsEmpty(rim.Diameter) {
		return r, errors.New("OnePiece Field Empty")
	}
	r.Diameter, err = db.ConvertStringToNullInt32(rim.Diameter)
	if err != nil {
		return r, err
	}

	//Material
	if db.IsEmpty(rim.Material) {
		return r, errors.New("OnePiece Field Empty")
	}
	r.Material = db.ConvertStringToNullString(rim.Material)
	if err != nil {
		return r, err
	}

	// Height
	if db.IsEmpty(rim.Height) {
		return r, errors.New("OnePiece Field Empty")
	}
	r.Height = db.ConvertStringToNullString(rim.Height)
	return r, err
}
