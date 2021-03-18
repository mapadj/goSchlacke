package timespans

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mapadj/goSchlacke/internal/pkg/convert"
)

/*
Feldname	 	Typ 	Länge 	Beschreibung
schwacke_id		Num		8		Interner Zeitraumschlüssel
schwacke_code	Num		8		Interner Typschlüssel
valid_from		Date	10		Datum (TT.MM.JJJJ)
valid_to		Date	10		Datum (TT.MM.JJJJ). Wert ist 00.00.0000 falls gültig
*/

// ConvertLineTimespansV1 convert one line
func (container TimespansV1Container) ConvertAndValidate() (err error) {

	timespan := container.GetImportStruct().(*importTable)

	time0 := time.Now()
	// SchwackeId, Num, Len 8
	if convert.IsEmpty(timespan.SchwackeID) {
		return errors.New("SchwackeID Field Empty")
	}

	container.ConvertedAndValidated.SchwackeID, err = convert.ConvertStringToInt32(timespan.SchwackeID)
	if err != nil {
		return fmt.Errorf("SchwackeId: Conversion Error: %v", err)
	}

	time1 := time.Now()
	// SchwackeId, Num, Len 8
	if convert.IsEmpty(timespan.SchwackeCode) {
		return errors.New("SchwackeCode Field Empty")
	}

	container.ConvertedAndValidated.SchwackeID, err = convert.ConvertStringToInt32(timespan.SchwackeCode)
	if err != nil {
		return fmt.Errorf("SchwackeCode: Conversion Error: %v", err)
	}

	time2 := time.Now()
	// valid_from, Len 10, Datum (TT.MM.JJJJ)
	if convert.IsEmpty(timespan.ValidFrom) {
		return errors.New("ValidFrom Field Empty")
	}

	container.ConvertedAndValidated.ValidFrom, err = convert.ConvertStringToNullDateStandardFormatWithAllZeroCheck(timespan.ValidFrom)
	if err != nil {
		println("VALUES: ", timespan.ValidFrom, timespan.ValidTo)
		return fmt.Errorf("ValidFrom: Conversion Error: %v", err)
	}

	time3 := time.Now()
	// valid_to, Len 10, Datum (TT.MM.JJJJ)
	if convert.IsEmpty(timespan.ValidTo) {

		return errors.New("ValidTo Field Empty")
	}

	container.ConvertedAndValidated.ValidUntil, err = convert.ConvertStringToNullDateStandardFormatWithAllZeroCheck(timespan.ValidTo)
	if err != nil {
		println("VALUES: ", timespan.ValidFrom, timespan.ValidTo)
		return fmt.Errorf("ValidTo: Conversion Error: %v", err)
	}

	time4 := time.Now()

	s1 := time1.Sub(time0).Microseconds()
	s2 := time2.Sub(time1).Microseconds()
	s3 := time3.Sub(time2).Microseconds()
	s4 := time4.Sub(time3).Microseconds()
	if s1+s2+s3+s4 > 0 {
		log.Printf("TIMES: \t%d\t%d\t%d\t%d DATA: %+v\n", s1, s2, s3, s4, timespan)
	}

	// Return Result
	return nil
}
