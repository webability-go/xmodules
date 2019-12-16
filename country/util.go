package country

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["country_country"] = country_country()
	sitecontext.Tables["country_country"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["country_country"].SetLanguage(language.English)
}

func buildCache(sitecontext *context.Context) {

	// Loads all data in XCache
	countries, _ := sitecontext.Tables["country_country"].SelectAll()

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

func SynchronizeDatabase(sitecontext *context.Context) {

	num1, err1 := sitecontext.Tables["country_country"].Count(nil)
	if err1 != nil || num1 == 0 {
		sitecontext.Logs["main"].Println("The table country_country was created (again)")
		sitecontext.Tables["country_country"].Synchronize()
	} else {
		sitecontext.Logs["main"].Println("The table country_country was not created because it contains data")
	}

	// fill countries and translations

}
