package material

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xmodules/base"
)

func buildTables(ctx *base.Datasource) {

	ctx.SetTable("material_material", materialMaterial())
	ctx.GetTable("material_material").SetBase(ctx.GetDatabase())
	ctx.GetTable("material_material").SetLanguage(language.Spanish)
}

func createCache(ds applications.Datasource) []string {

	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()
		ds.SetCache("materiales:"+canonical, xcore.NewXCache("materiales:"+canonical, 0, 0))
	}
	return []string{}
}

func buildCache(ds applications.Datasource) []string {

	// Lets protect us for race condition since map[] of Tables and XCaches are not thread safe
	material_material := ds.GetTable("material_material")
	caches := map[string]*xcore.XCache{}
	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()
		caches["materiales:"+canonical] = ds.GetCache("materiales:" + canonical)
	}

	// Loads all data in XCache
	materiales, _ := material_material.SelectAll()
	if materiales == nil {
		return []string{"No hay material en la base de datos"}
	}
	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()

		all := []int{}
		for _, m := range *materiales {
			// creates structure on language
			str := CreateStructureMaterialByData(ds, m.Clone(), lang)
			key, _ := m.GetInt("clave")
			all = append(all, key)
			caches["materiales:"+canonical].Set(strconv.Itoa(key), str)
		}
		caches["materiales:"+canonical].Set("all", all)
	}

	return []string{}
}
