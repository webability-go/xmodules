package metric

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["metric_unit"] = metricUnit()
	sitecontext.Tables["metric_unit"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["metric_unit"].SetLanguage(language.English)
}

func buildCache(sitecontext *context.Context) {
	// Loads all data in XCache
	metrics, _ := sitecontext.Tables["metric_unit"].SelectAll()
	if metrics == nil {
		return
	}

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["metric:"+canonical] = xcore.NewXCache("metric:"+canonical, 0, 0)

		for _, m := range *metrics {
			// creates structure on language
			str := CreateStructureMetricByData(sitecontext, m.Clone(), lang)
			clave, _ := m.GetInt("key")
			sitecontext.Caches["metric:"+canonical].Set(strconv.Itoa(clave), str)
		}
	}
}
