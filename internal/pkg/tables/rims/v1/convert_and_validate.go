package rims

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mapadj/goSchlacke/internal/pkg/convert"
)

// ConvertAndValidate Convert Function
func (rc RimsV1Container) ConvertAndValidate() (err error) {

	// Init vars
	rim := rc.RawStruct

	// CODE
	if convert.IsEmpty(rim.Code) {
		return errors.New("Code Field Empty")
	}
	rc.ConvertedAndValidated.Code, err = convert.ConvertStringToInt32(rim.Code)
	if err != nil {
		return err
	}

	// WIDTH
	if convert.IsEmpty(rim.Width) {
		errors.New("Width Field Empty")
	}
	w := strings.TrimSpace(rim.Width)
	_, err = strconv.ParseFloat(w, 64)
	if err != nil {
		return fmt.Errorf("Width: Float Conversion Error: %v", err)
	}
	rc.ConvertedAndValidated.Width = convert.ConvertStringToNullString(w)
	if err != nil {
		return err
	}

	// OnePiece
	if convert.IsEmpty(rim.OnePiece) {
		return errors.New("OnePiece Field Empty")
	}
	rc.ConvertedAndValidated.OnePiece, err = convert.ConvertXxToNullBool(rim.OnePiece)
	if err != nil {
		return err
	}

	// Diameter
	if convert.IsEmpty(rim.Diameter) {
		return errors.New("OnePiece Field Empty")
	}
	rc.ConvertedAndValidated.Diameter, err = convert.ConvertStringToNullInt32(rim.Diameter)
	if err != nil {
		return err
	}

	//Material
	if convert.IsEmpty(rim.Material) {
		return errors.New("OnePiece Field Empty")
	}
	rc.ConvertedAndValidated.Material = convert.ConvertStringToNullString(rim.Material)
	if err != nil {
		return err
	}

	// Height
	if convert.IsEmpty(rim.Height) {
		return errors.New("OnePiece Field Empty")
	}
	rc.ConvertedAndValidated.Height = convert.ConvertStringToNullString(rim.Height)

	//

	return nil
}
