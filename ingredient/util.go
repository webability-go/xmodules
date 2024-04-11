package ingredient

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xcore/v2"
)

func buildTables(ds applications.Datasource) {

	ds.SetTable("ingredient_aisle", ingredientAisle())
	ds.GetTable("ingredient_aisle").SetBase(ds.GetDatabase())
	ds.GetTable("ingredient_aisle").SetLanguage(language.Spanish)

	ds.SetTable("ingredient_ingredient", ingredientIngredient())
	ds.GetTable("ingredient_ingredient").SetBase(ds.GetDatabase())
	ds.GetTable("ingredient_ingredient").SetLanguage(language.Spanish)
}

func createCache(ds applications.Datasource) []string {

	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()
		ds.SetCache("ingredient:pasillos:"+canonical, xcore.NewXCache("ingredient:pasillos:"+canonical, 0, 0))
		ds.SetCache("ingredient:ingredientes:"+canonical, xcore.NewXCache("ingredient:ingredientes:"+canonical, 0, 0))
	}
	return []string{}
}

func buildCache(ds applications.Datasource) []string {

	// Lets protect us for race condition since map[] of Tables and XCaches are not thread safe
	ingredient_aisle := ds.GetTable("ingredient_aisle")
	ingredient_ingredient := ds.GetTable("ingredient_ingredient")
	caches := map[string]*xcore.XCache{}
	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()
		caches["ingredient:pasillos:"+canonical] = ds.GetCache("ingredient:pasillos:" + canonical)
		caches["ingredient:ingredientes:"+canonical] = ds.GetCache("ingredient:ingredientes:" + canonical)
	}

	// Loads all data in XCache
	pasillos, _ := ingredient_aisle.SelectAll()

	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()

		all := []int{}
		if pasillos != nil {
			for _, m := range *pasillos {
				// creates structure on language
				str := CreateStructurePasilloByData(ds, m.Clone(), lang)
				key, _ := m.GetInt("clave")
				all = append(all, key)
				caches["ingredient:pasillos:"+canonical].Set(strconv.Itoa(key), str)
			}
		}
		caches["ingredient:pasillos:"+canonical].Set("all", all)
	}

	// Loads all data in XCache
	ingredients, _ := ingredient_ingredient.SelectAll()

	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()

		if ingredients != nil {
			for _, m := range *ingredients {
				// creates structure on language
				str := CreateStructureIngredientByData(ds, m.Clone(), lang)
				key, _ := m.GetInt("clave")
				caches["ingredient:ingredientes:"+canonical].Set(strconv.Itoa(key), str)
			}
		}
	}

	return []string{}
}
