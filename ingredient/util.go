package ingredient

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xmodules/base"
)

func buildTables(ctx *base.Datasource) {

	ctx.SetTable("ingredient_aisle", ingredientAisle())
	ctx.GetTable("ingredient_aisle").SetBase(ctx.GetDatabase())
	ctx.GetTable("ingredient_aisle").SetLanguage(language.Spanish)

	ctx.SetTable("ingredient_ingredient", ingredientIngredient())
	ctx.GetTable("ingredient_ingredient").SetBase(ctx.GetDatabase())
	ctx.GetTable("ingredient_ingredient").SetLanguage(language.Spanish)
}

func createCache(ctx *base.Datasource) []string {

	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()
		ctx.SetCache("ingredient:pasillos:"+canonical, xcore.NewXCache("ingredient:pasillos:"+canonical, 0, 0))
		ctx.SetCache("ingredient:ingredientes:"+canonical, xcore.NewXCache("ingredient:ingredientes:"+canonical, 0, 0))
	}
	return []string{}
}

func buildCache(ctx *base.Datasource) []string {

	// Lets protect us for race condition since map[] of Tables and XCaches are not thread safe
	ingredient_aisle := ctx.GetTable("ingredient_aisle")
	ingredient_ingredient := ctx.GetTable("ingredient_ingredient")
	caches := map[string]*xcore.XCache{}
	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()
		caches["ingredient:pasillos:"+canonical] = ctx.GetCache("ingredient:pasillos:" + canonical)
		caches["ingredient:ingredientes:"+canonical] = ctx.GetCache("ingredient:ingredientes:" + canonical)
	}

	// Loads all data in XCache
	pasillos, _ := ingredient_aisle.SelectAll()

	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()

		all := []int{}
		if pasillos != nil {
			for _, m := range *pasillos {
				// creates structure on language
				str := CreateStructurePasilloByData(ctx, m.Clone(), lang)
				key, _ := m.GetInt("clave")
				all = append(all, key)
				caches["ingredient:pasillos:"+canonical].Set(strconv.Itoa(key), str)
			}
		}
		caches["ingredient:pasillos:"+canonical].Set("all", all)
	}

	// Loads all data in XCache
	ingredients, _ := ingredient_ingredient.SelectAll()

	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()

		if ingredients != nil {
			for _, m := range *ingredients {
				// creates structure on language
				str := CreateStructureIngredientByData(ctx, m.Clone(), lang)
				key, _ := m.GetInt("clave")
				caches["ingredient:ingredientes:"+canonical].Set(strconv.Itoa(key), str)
			}
		}
	}

	return []string{}
}
