// Package country contains the list of all countries in the world with some data
// Necesita USDA y METRICS para funcionar
package country

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID         = "country"
	VERSION          = "1.0.0"
	TRANSLATIONTHEME = "country"
)

// InitCountry is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func InitModule(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	buildCache(sitecontext)
	sitecontext.Modules[MODULEID] = VERSION

	return nil
}

func SynchronizeModule(sitecontext *context.Context) []string {

	translation.AddTheme(sitecontext, TRANSLATIONTHEME, "Countries", translation.SOURCETABLE, "", "name")

	messages := []string{}
	messages = append(messages, "Analysing country_country table.")
	num, err := sitecontext.Tables["country_country"].Count(nil)
	if err != nil || num == 0 {
		err1 := sitecontext.Tables["country_country"].Synchronize()
		if err1 != nil {
			messages = append(messages, "The table country_country was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table country_country was created (again)")
		}
	} else {
		messages = append(messages, "The table country_country was not created because it contains data.")
	}

	// fill countries and translations

	return messages
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetCountry(sitecontext *context.Context, key string, lang language.Tag) *StructureCountry {

	canonical := lang.String()

	data, _ := sitecontext.Caches["country:countries:"+canonical].Get(key)
	if data == nil {
		sm := CreateStructureCountryByKey(sitecontext, key, lang)
		if sm == nil {
			sitecontext.Logs["graph"].Println("xmodules::country::GetGountry: there is no country created:", key, lang)
			return nil
		}
		sitecontext.Caches["country:countries:"+canonical].Set(key, sm)
		return sm.(*StructureCountry)
	}
	return data.(*StructureCountry)
}
