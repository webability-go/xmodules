package country

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["country_country"] = countryCountry()
	sitecontext.Tables["country_country"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["country_country"].SetLanguage(language.English)
}

func buildCache(sitecontext *context.Context) {

	// Loads all data in XCache
	countries, _ := sitecontext.Tables["country_country"].SelectAll()
	if countries == nil {
		return
	}

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["country:countries:"+canonical] = xcore.NewXCache("country:countries:"+canonical, 0, 0)

		all := []string{}
		for _, m := range *countries {
			// creates structure on language
			str := CreateStructureCountryByData(sitecontext, m.Clone(), lang)
			key, _ := m.GetString("clave")
			all = append(all, key)
			sitecontext.Caches["country:countries:"+canonical].Set(key, str)
		}
		sitecontext.Caches["country:countries:"+canonical].Set("all", all)
	}
}
