// Package country contains the list of all countries in the world with some data
// Necesita USDA y METRICS para funcionar
package country

import (
	"github.com/webability-go/xamboo/applications"
	"golang.org/x/text/language"
)

func createTables(ds applications.Datasource) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {
		messages = append(messages, "Analysing "+tbl+" table.")
		num, err := ds.GetTable(tbl).Count(nil)
		if err != nil || num == 0 {
			err1 := ds.GetTable(tbl).Synchronize()
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
func GetCountry(ds applications.Datasource, key string, lang language.Tag) *StructureCountry {

	canonical := lang.String()

	data, _ := ds.GetCache("country:countries:" + canonical).Get(key)
	if data == nil {
		sm := CreateStructureCountryByKey(ds, key, lang)
		if sm == nil {
			ds.Log("graph", "xmodules::country::GetGountry: there is no country created:", key, lang)
			return nil
		}
		ds.GetCache("country:countries:"+canonical).Set(key, sm)
		return sm.(*StructureCountry)
	}
	return data.(*StructureCountry)
}
