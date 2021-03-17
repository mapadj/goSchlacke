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
type importTable struct {
	tables.ImportStruct
	SchwackeID   string `fixed:"1,8"`
	SchwackeCode string `fixed:"9,16"`
	ValidFrom    string `fixed:"17,26"`
	ValidTo      string `fixed:"27,36"`
}

// Container
type TimespansV1Container struct {
	importTable           *importTable
	ConvertedAndValidated *db.UpsertTimespansV1Params
}

func (c TimespansV1Container) GetImportStruct() tables.ImportStruct {
	return c.importTable
}

func (c TimespansV1Container) GetUpsertStruct() interface{} {
	return c.ConvertedAndValidated
}

// Handler
type TimespansV1Handler struct {
}

func NewHandler() tables.ImportHandler {
	return &TimespansV1Handler{}
}

func (handler TimespansV1Handler) NewContainer() tables.ImportContainer {
	return &TimespansV1Container{
		importTable:           &importTable{},
		ConvertedAndValidated: &db.UpsertTimespansV1Params{},
	}
}
