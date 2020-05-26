// Package country contains the list of all countries in the world with some data
// Necesita USDA y METRICS para funcionar
package country

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/base"
)

func createTables(ctx *base.Datasource) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {
		messages = append(messages, "Analysing "+tbl+" table.")
		num, err := ctx.GetTable(tbl).Count(nil)
		if err != nil || num == 0 {
			err1 := ctx.GetTable(tbl).Synchronize()
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
func GetCountry(ctx *base.Datasource, key string, lang language.Tag) *StructureCountry {

	canonical := lang.String()

	data, _ := ctx.GetCache("country:countries:" + canonical).Get(key)
	if data == nil {
		sm := CreateStructureCountryByKey(ctx, key, lang)
		if sm == nil {
			ctx.Log("graph", "xmodules::country::GetGountry: there is no country created:", key, lang)
			return nil
		}
		ctx.GetCache("country:countries:"+canonical).Set(key, sm)
		return sm.(*StructureCountry)
	}
	return data.(*StructureCountry)
}
