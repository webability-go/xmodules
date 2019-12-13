package ingredient

import (
  "strconv"

  "golang.org/x/text/language"

  "github.com/webability-go/xcore"

  "xmodules/context"
//  "xmodules/translation"
)

func buildTables(sitecontext *context.Context, databasename string) {

    sitecontext.Tables["kl_ingredientepasillo"] = kl_ingredientepasillo()
    sitecontext.Tables["kl_ingredientepasillo"].SetBase(sitecontext.Databases[databasename])
    sitecontext.Tables["kl_ingredientepasillo"].SetLanguage(language.Spanish)

    sitecontext.Tables["kl_ingrediente"] = kl_ingrediente()
    sitecontext.Tables["kl_ingrediente"].SetBase(sitecontext.Databases[databasename])
    sitecontext.Tables["kl_ingrediente"].SetLanguage(language.Spanish)
}

func buildCache(sitecontext *context.Context) {

  // Loads all data in XCache
  pasillos, _ := sitecontext.Tables["kl_ingredientepasillo"].SelectAll()

  for _, lang := range sitecontext.Languages {
	  canonical := lang.String()
    sitecontext.Caches["ingredient:pasillos:" + canonical] = xcore.NewXCache("ingredient:pasillos:" + canonical, 0, 0)

    all := []int{}
		for _, m := range *pasillos {
      // creates structure on language
			str := CreateStructurePasilloByData(sitecontext, m.Clone(), lang)
			key, _ := m.GetInt("clave")
      all = append(all, key)
			sitecontext.Caches["ingredient:pasillos:" + canonical].Set(strconv.Itoa(key), str)
		}
    sitecontext.Caches["ingredient:pasillos:" + canonical].Set("all", all)
	}

  // Loads all data in XCache
  ingredients, _ := sitecontext.Tables["kl_ingrediente"].SelectAll()

  for _, lang := range sitecontext.Languages {
	  canonical := lang.String()
    sitecontext.Caches["ingredient:ingredientes:" + canonical] = xcore.NewXCache("ingredient:ingredientes:" + canonical, 0, 0)

		for _, m := range *ingredients {
      // creates structure on language
			str := CreateStructureIngredientByData(sitecontext, m.Clone(), lang)
			key, _ := m.GetInt("clave")
			sitecontext.Caches["ingredient:ingredientes:" + canonical].Set(strconv.Itoa(key), str)
		}
	}

}

func SynchronizeDatabase(sitecontext *context.Context) {

  num1, err1 := sitecontext.Tables["kl_ingredientepasillo"].Count(nil)
  num2, err2 := sitecontext.Tables["kl_ingrediente"].Count(nil)
  if (err1 != nil && err2 != nil) || (num1 == 0 && num2 == 0) {
    sitecontext.Logs["main"].Println("The tables kl_ingrediente* were created (again)")
    sitecontext.Tables["kl_ingredientepasillo"].Synchronize()
    sitecontext.Tables["kl_ingrediente"].Synchronize()
  } else {
    sitecontext.Logs["main"].Println("The tables kl_ingrediente* were not created because they contain data")
  }
}
