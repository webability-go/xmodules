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
	VERSION          = "2.0.0"
	TRANSLATIONTHEME = "country"
)

// InitCountry is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func InitModule(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	createCache(sitecontext)
	sitecontext.SetModule(MODULEID, VERSION)

	go buildCache(sitecontext)

	return nil
}

func SynchronizeModule(sitecontext *context.Context, filespath string) []string {

	messages := []string{}

	// Needed modules: context and translation
	vc := context.ModuleInstalledVersion(sitecontext, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/country.")
		return messages
	}
	vc = context.ModuleInstalledVersion(sitecontext, "translation")
	if vc == "" {
		messages = append(messages, "xmodules/translation need to be installed before installing xmodules/country.")
		return messages
	}

	translation.AddTheme(sitecontext, TRANSLATIONTHEME, "Countries", translation.SOURCETABLE, "", "name")

	// create tables
	messages = append(messages, createTables(sitecontext)...)

	// fill countries and translations
	messages = append(messages, loadTables(sitecontext, filespath)...)

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := context.AddModule(sitecontext, MODULEID, "List of official countries and ISO codes", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages
}

func createTables(sitecontext *context.Context) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {
		messages = append(messages, "Analysing "+tbl+" table.")
		num, err := sitecontext.GetTable(tbl).Count(nil)
		if err != nil || num == 0 {
			err1 := sitecontext.GetTable(tbl).Synchronize()
			if err1 != nil {
				messages = append(messages, "The table "+tbl+" was not created: "+err1.Error())
			} else {
				messages = append(messages, "The table "+tbl+" was created (again)")
			}
		} else {
			messages = append(messages, "The table "+tbl+" was not created because it contains data.")
		}
	}

	return messages
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetCountry(sitecontext *context.Context, key string, lang language.Tag) *StructureCountry {

	canonical := lang.String()

	data, _ := sitecontext.GetCache("country:countries:" + canonical).Get(key)
	if data == nil {
		sm := CreateStructureCountryByKey(sitecontext, key, lang)
		if sm == nil {
			sitecontext.Log("graph", "xmodules::country::GetGountry: there is no country created:", key, lang)
			return nil
		}
		sitecontext.GetCache("country:countries:"+canonical).Set(key, sm)
		return sm.(*StructureCountry)
	}
	return data.(*StructureCountry)
}
