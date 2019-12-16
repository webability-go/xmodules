// Package country contains the list of all countries in the world with some data
// Necesita USDA y METRICS para funcionar
package country

import (
	"fmt"
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/metrics"
)

// InitCountry is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func InitCountry(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	buildCache(sitecontext)

	return nil
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetCountry(sitecontext *context.Context, key string, lang language.Tag) *StructureCountry {

	canonical := lang.String()

	data, _ := sitecontext.Caches["country:countries:"+canonical].Get(clave)
	if data == nil {
		sm := CreateStructureCountryByKey(sitecontext, clave, lang)
		if sm == nil {
			sitecontext.Logs["graph"].Println("xmodules::country::GetGountry: there is no country created:", clave, lang)
			return nil
		}
		sitecontext.Caches["country:countries:"+canonical].Set(clave, sm)
		return sm.(*StructureCountry)
	}
	return data.(*StructureCountry)
}
