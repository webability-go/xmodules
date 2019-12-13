package metrics

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"

	"xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["metric_unit"] = kl_medida()
	sitecontext.Tables["metric_unit"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["metric_unit"].SetLanguage(language.English)
}

func buildCache(sitecontext *context.Context) {
	// Loads all data in XCache
	metrics, _ := sitecontext.Tables["metric_unit"].SelectAll()

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["metrics:"+canonical] = xcore.NewXCache("metrics:"+canonical, 0, 0)

		for _, m := range *metrics {
			// creates structure on language
			str := CreateStructureMetricByData(sitecontext, m.Clone(), lang)
			clave, _ := m.GetInt("key")
			sitecontext.Caches["metrics:"+canonical].Set(strconv.Itoa(clave), str)
		}
	}
}

func SynchronizeDatabase(sitecontext *context.Context) {

	num, err := sitecontext.Tables["metric_unit"].Count(nil)
	if err != nil || num == 0 {
		sitecontext.Logs["main"].Println("The table metric_unit was created (again)")
		sitecontext.Tables["metric_unit"].Synchronize()
	} else {
		sitecontext.Logs["main"].Println("The table metric_unit was not created because it contains data")
	}
}
