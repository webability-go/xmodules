package metric

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xmodules/base"
)

func buildTables(ctx *base.Datasource) {

	ctx.SetTable("metric_unit", metricUnit())
	ctx.GetTable("metric_unit").SetBase(ctx.GetDatabase())
	ctx.GetTable("metric_unit").SetLanguage(language.English)
}

func createCache(ctx *base.Datasource) []string {

	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()
		ctx.SetCache("metric:"+canonical, xcore.NewXCache("metric:"+canonical, 0, 0))
	}
	return []string{}
}

func buildCache(ctx *base.Datasource) []string {

	// Lets protect us for race condition since map[] of Tables and XCaches are not thread safe
	metric_unit := ctx.GetTable("metric_unit")
	caches := map[string]*xcore.XCache{}
	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()
		caches["metric:"+canonical] = ctx.GetCache("metric:" + canonical)
	}

	// Loads all data in XCache
	metrics, _ := metric_unit.SelectAll()
	if metrics == nil {
		return []string{"No hay metricas en la base de datos"}
	}

	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()

		for _, m := range *metrics {
			// creates structure on language
			str := CreateStructureMetricByData(ctx, m.Clone(), lang)
			clave, _ := m.GetInt("key")
			caches["metric:"+canonical].Set(strconv.Itoa(clave), str)
		}
	}

	return []string{}
}
