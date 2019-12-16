package context

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["context_module"] = contextModule()
	sitecontext.Tables["context_module"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["context_module"].SetLanguage(language.English)
}

func buildCache(sitecontext *Context) {

	// Loads all data in XCache
	modules, _ := sitecontext.Tables["context_module"].SelectAll()

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["context:modules:"+canonical] = xcore.NewXCache("context:modules:"+canonical, 0, 0)

		all := []string{}
		for _, m := range *modules {
			// creates structure on language
			str := CreateStructureModuleByData(sitecontext, m.Clone(), lang)
			key, _ := m.GetString("key")
			all = append(all, key)
			sitecontext.Caches["context:modules:"+canonical].Set(key, str)
		}
		sitecontext.Caches["context:modules:"+canonical].Set("all", all)
	}
}

func SynchronizeDatabase(sitecontext *Context) {

	num1, err1 := sitecontext.Tables["context_module"].Count(nil)
	if err1 != nil || num1 == 0 {
		sitecontext.Logs["main"].Println("The table context_module was created (again)")
		sitecontext.Tables["context_module"].Synchronize()
	} else {
		sitecontext.Logs["main"].Println("The table context_module was not created because it contains data")
	}

	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	sitecontext.Tables["context_module"].Upsert(&xdominion.XRecord{
		"key":     "context",
		"name":    "Contexts for Xamboo",
		"version": "1.0.0",
	})
}
