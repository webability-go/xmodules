package metric

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xcore/v2"
)

func buildTables(ds applications.Datasource) {

	ds.SetTable("metric_unit", metricUnit())
	ds.GetTable("metric_unit").SetBase(ds.GetDatabase())
	ds.GetTable("metric_unit").SetLanguage(language.English)
}

func createCache(ds applications.Datasource) []string {

	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()
		ds.SetCache("metric:"+canonical, xcore.NewXCache("metric:"+canonical, 0, 0))
	}
	return []string{}
}

func buildCache(ds applications.Datasource) []string {

	// Lets protect us for race condition since map[] of Tables and XCaches are not thread safe
	metric_unit := ds.GetTable("metric_unit")
	caches := map[string]*xcore.XCache{}
	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()
		caches["metric:"+canonical] = ds.GetCache("metric:" + canonical)
	}

	// Loads all data in XCache
	metrics, _ := metric_unit.SelectAll()
	if metrics == nil {
		return []string{"No hay metricas en la base de datos"}
	}

	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()

		for _, m := range *metrics {
			// creates structure on language
			str := CreateStructureMetricByData(ds, m.Clone(), lang)
			clave, _ := m.GetInt("key")
			caches["metric:"+canonical].Set(strconv.Itoa(clave), str)
		}
	}

	return []string{}
}
