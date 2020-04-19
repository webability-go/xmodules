package metric

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.SetTable("metric_unit", metricUnit())
	sitecontext.GetTable("metric_unit").SetBase(sitecontext.GetDatabase(databasename))
	sitecontext.GetTable("metric_unit").SetLanguage(language.English)
}

func createCache(sitecontext *context.Context) []string {

	for _, lang := range sitecontext.GetLanguages() {
		canonical := lang.String()
		sitecontext.SetCache("metric:"+canonical, xcore.NewXCache("metric:"+canonical, 0, 0))
	}
	return []string{}
}

func buildCache(sitecontext *context.Context) []string {

	// Lets protect us for race condition since map[] of Tables and XCaches are not thread safe
	metric_unit := sitecontext.GetTable("metric_unit")
	caches := map[string]*xcore.XCache{}
	for _, lang := range sitecontext.GetLanguages() {
		canonical := lang.String()
		caches["metric:"+canonical] = sitecontext.GetCache("metric:" + canonical)
	}

	// Loads all data in XCache
	metrics, _ := metric_unit.SelectAll()
	if metrics == nil {
		return []string{"No hay metricas en la base de datos"}
	}

	for _, lang := range sitecontext.GetLanguages() {
		canonical := lang.String()

		for _, m := range *metrics {
			// creates structure on language
			str := CreateStructureMetricByData(sitecontext, m.Clone(), lang)
			clave, _ := m.GetInt("key")
			caches["metric:"+canonical].Set(strconv.Itoa(clave), str)
		}
	}

	return []string{}
}
