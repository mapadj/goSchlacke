package SchlackeImporter

import (
	"reflect"

	"github.com/mapadj/goSchlacke/internal/pkg/tables"
	"github.com/mapadj/goSchlacke/internal/pkg/tables/rims/v1"
	"github.com/mapadj/goSchlacke/internal/pkg/tables/timespans/v1"
)

// TableRegistry contains all different TableTypes
var TableRegistry = map[string]reflect.Type{
	"RimsV1": reflect.TypeOf((*rims.RimsV1Container)(nil)).Elem(),
}

// FactoryMap contains all different Importable Factories
var ContainerFactoryMap = map[string]func() tables.ImportHandler{
	"RimsV1":      rims.NewHandler,
	"TimespansV1": timespans.NewHandler,
}
