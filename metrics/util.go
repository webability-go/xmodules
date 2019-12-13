package metrics

import (
  "strconv"

  "golang.org/x/text/language"

  "github.com/webability-go/xcore"

  "xmodules/context"
//  "xmodules/translation"
)

func buildTables(sitecontext *context.Context, databasename string) {

  sitecontext.Tables["kl_medida"] = kl_medida()
  sitecontext.Tables["kl_medida"].SetBase(sitecontext.Databases[databasename])
  sitecontext.Tables["kl_medida"].SetLanguage(language.Spanish)
}

func buildCache(sitecontext *context.Context) {
  // Loads all data in XCache
  metrics, _ := sitecontext.Tables["kl_medida"].SelectAll()

  for _, lang := range sitecontext.Languages {
	  canonical := lang.String()
    sitecontext.Caches["metrics:" + canonical] = xcore.NewXCache("metrics:" + canonical, 0, 0)

		for _, m := range *metrics {
      // creates structure on language
			str := CreateStructureMetricByData(sitecontext, m.Clone(), lang)
			clave, _ := m.GetInt("clave")
			sitecontext.Caches["metrics:" + canonical].Set(strconv.Itoa(clave), str)
		}
	}
}

func SynchronizeDatabase(sitecontext *context.Context) {

  num, err := sitecontext.Tables["kl_medida"].Count(nil)
  if err != nil || num == 0 {
    sitecontext.Logs["main"].Println("The table kl_medida was created (again)")
    sitecontext.Tables["kl_medida"].Synchronize()
  } else {
    sitecontext.Logs["main"].Println("The table kl_medida was not created because it contains data")
  }
}
