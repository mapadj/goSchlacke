package timespans

import (
	db "github.com/mapadj/goSchlacke/internal/pkg/db/sqlc"
	"github.com/mapadj/goSchlacke/internal/pkg/tables"
)

/*
Feldname 		Typ 	Länge 	Beschreibung

schwacke_id 	Num 	8
schwacke_code 	Num		8
valid_from		Char	10
valid_to		Bool	10

*/
// TimespansV1LineLength Length of one Line
const TimespansV1LineLength = 36

// ImportTable define the format
type ImportTable struct {
	SchwackeID   string `fixed:"1,8"`
	SchwackeCode string `fixed:"9,16"`
	ValidFrom    string `fixed:"17,26"`
	ValidTo      string `fixed:"27,36"`
}

type TimespansV1Handler struct {
}

func NewHandler() tables.ImportHandler {
	return TimespansV1Handler{}
}

type TimespansV1Container struct {
	RawStruct             ImportTable
	ConvertedAndValidated db.UpsertTimespansV1Params
}

func (handler TimespansV1Handler) NewContainer() tables.Importable {
	return TimespansV1Container{
		// RawStruct: ImportTable{},
	}
}
