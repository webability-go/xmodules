package ingredient

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.SetTable("ingredient_aisle", ingredientAisle())
	sitecontext.GetTable("ingredient_aisle").SetBase(sitecontext.GetDatabase(databasename))
	sitecontext.GetTable("ingredient_aisle").SetLanguage(language.Spanish)

	sitecontext.SetTable("ingredient_ingredient", ingredientIngredient())
	sitecontext.GetTable("ingredient_ingredient").SetBase(sitecontext.GetDatabase(databasename))
	sitecontext.GetTable("ingredient_ingredient").SetLanguage(language.Spanish)
}

func createCache(sitecontext *context.Context) []string {

	for _, lang := range sitecontext.GetLanguages() {
		canonical := lang.String()
		sitecontext.SetCache("ingredient:pasillos:"+canonical, xcore.NewXCache("ingredient:pasillos:"+canonical, 0, 0))
		sitecontext.SetCache("ingredient:ingredientes:"+canonical, xcore.NewXCache("ingredient:ingredientes:"+canonical, 0, 0))
	}
	return []string{}
}

func buildCache(sitecontext *context.Context) []string {

	// Lets protect us for race condition since map[] of Tables and XCaches are not thread safe
	ingredient_aisle := sitecontext.GetTable("ingredient_aisle")
	ingredient_ingredient := sitecontext.GetTable("ingredient_ingredient")
	caches := map[string]*xcore.XCache{}
	for _, lang := range sitecontext.GetLanguages() {
		canonical := lang.String()
		caches["ingredient:pasillos:"+canonical] = sitecontext.GetCache("ingredient:pasillos:" + canonical)
		caches["ingredient:ingredientes:"+canonical] = sitecontext.GetCache("ingredient:ingredientes:" + canonical)
	}

	// Loads all data in XCache
	pasillos, _ := ingredient_aisle.SelectAll()

	for _, lang := range sitecontext.GetLanguages() {
		canonical := lang.String()

		all := []int{}
		if pasillos != nil {
			for _, m := range *pasillos {
				// creates structure on language
				str := CreateStructurePasilloByData(sitecontext, m.Clone(), lang)
				key, _ := m.GetInt("clave")
				all = append(all, key)
				caches["ingredient:pasillos:"+canonical].Set(strconv.Itoa(key), str)
			}
		}
		caches["ingredient:pasillos:"+canonical].Set("all", all)
	}

	// Loads all data in XCache
	ingredients, _ := ingredient_ingredient.SelectAll()

	for _, lang := range sitecontext.GetLanguages() {
		canonical := lang.String()

		if ingredients != nil {
			for _, m := range *ingredients {
				// creates structure on language
				str := CreateStructureIngredientByData(sitecontext, m.Clone(), lang)
				key, _ := m.GetInt("clave")
				caches["ingredient:ingredientes:"+canonical].Set(strconv.Itoa(key), str)
			}
		}
	}

	return []string{}
}
