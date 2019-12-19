package ingredient

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["ingredient_aisle"] = ingredientAisle()
	sitecontext.Tables["ingredient_aisle"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["ingredient_aisle"].SetLanguage(language.Spanish)

	sitecontext.Tables["ingredient_ingredient"] = ingredientIngredient()
	sitecontext.Tables["ingredient_ingredient"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["ingredient_ingredient"].SetLanguage(language.Spanish)
}

func buildCache(sitecontext *context.Context) {

	// Loads all data in XCache
	pasillos, _ := sitecontext.Tables["ingredient_aisle"].SelectAll()

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["ingredient:pasillos:"+canonical] = xcore.NewXCache("ingredient:pasillos:"+canonical, 0, 0)

		all := []int{}
		if pasillos != nil {
			for _, m := range *pasillos {
				// creates structure on language
				str := CreateStructurePasilloByData(sitecontext, m.Clone(), lang)
				key, _ := m.GetInt("clave")
				all = append(all, key)
				sitecontext.Caches["ingredient:pasillos:"+canonical].Set(strconv.Itoa(key), str)
			}
		}
		sitecontext.Caches["ingredient:pasillos:"+canonical].Set("all", all)
	}

	// Loads all data in XCache
	ingredients, _ := sitecontext.Tables["ingredient_ingredient"].SelectAll()

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["ingredient:ingredientes:"+canonical] = xcore.NewXCache("ingredient:ingredientes:"+canonical, 0, 0)

		if ingredients != nil {
			for _, m := range *ingredients {
				// creates structure on language
				str := CreateStructureIngredientByData(sitecontext, m.Clone(), lang)
				key, _ := m.GetInt("clave")
				sitecontext.Caches["ingredient:ingredientes:"+canonical].Set(strconv.Itoa(key), str)
			}
		}
	}

}
