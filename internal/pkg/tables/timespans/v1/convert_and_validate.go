package timespans

import (
	"errors"
	"fmt"

	"github.com/mapadj/goSchlacke/internal/pkg/convert"
)

/*
Feldname	 	Typ 	Länge 	Beschreibung
schwacke_id		Num		8		Interner Zeitraumschlüssel
schwacke_code	Num		8		Interner Typschlüssel
valid_from		Date	10		Datum (TT.MM.JJJJ)
valid_to		Date	10		Datum (TT.MM.JJJJ). Wert ist 00.00.0000 falls gültig
*/

const date_layout string = "02.01.2006"

// ConvertLineTimespansV1 convert one line
func (container TimespansV1Container) ConvertAndValidate() (err error) {

	timespan := container.RawStruct

	// SchwackeId, Num, Len 8
	if convert.IsEmpty(timespan.SchwackeID) {
		return errors.New("SchwackeID Field Empty")
	}

	container.ConvertedAndValidated.SchwackeID, err = convert.ConvertStringToInt32(timespan.SchwackeID)
	if err != nil {
		return fmt.Errorf("SchwackeId: Conversion Error: %v", err)
	}

	// SchwackeId, Num, Len 8
	if convert.IsEmpty(timespan.SchwackeCode) {
		return errors.New("SchwackeCode Field Empty")
	}

	container.ConvertedAndValidated.SchwackeID, err = convert.ConvertStringToInt32(timespan.SchwackeCode)
	if err != nil {
		return fmt.Errorf("SchwackeCode: Conversion Error: %v", err)
	}

	// valid_from, Len 10, Datum (TT.MM.JJJJ)
	if convert.IsEmpty(timespan.ValidFrom) {
		return errors.New("ValidFrom Field Empty")
	}

	container.ConvertedAndValidated.ValidFrom, err = convert.ConvertStringToNullDate(timespan.ValidFrom, date_layout)
	if err != nil {
		return fmt.Errorf("ValidFrom: Conversion Error: %v", err)
	}

	// valid_to, Len 10, Datum (TT.MM.JJJJ)
	if convert.IsEmpty(timespan.ValidFrom) {
		return errors.New("ValidFrom Field Empty")
	}

	container.ConvertedAndValidated.ValidUntil, err = convert.ConvertStringToNullDate(timespan.ValidTo, date_layout)
	if err != nil {
		return fmt.Errorf("ValidTo: Conversion Error: %v", err)
	}

	// Return Result
	return nil
}
