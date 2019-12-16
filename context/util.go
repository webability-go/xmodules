package context

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore"
	"github.com/webability-go/xdominion"
)

func buildTables(sitecontext *Context, databasename string) {

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

// Analyze a context and gets back the main data
func GetContextStats(sitecontext *Context) *xcore.XDataset {

	subdata := xcore.XDataset{}
	subdata["languages"] = sitecontext.Languages
	subdata["modules"], _ = sitecontext.Config.Get("module")
	subdata["databases"] = sitecontext.Databases
	subdata["logs"] = sitecontext.Logs

	caches := []string{}
	for id := range sitecontext.Caches {
		caches = append(caches, id)
	}
	subdata["xcaches"] = caches

	tables := map[string]string{}
	for id, table := range sitecontext.Tables {
		db := table.Base.Database
		tables[id] = db
	}
	subdata["tables"] = tables

	subdata["config"] = buildConfigSet(sitecontext.Config)
	return &subdata
}

func buildConfigSet(config *xconfig.XConfig) xcore.XDataset {
	data := xcore.XDataset{}
	for id := range config.Parameters {
		d, _ := config.Get(id)
		if val, ok := d.(*xconfig.XConfig); ok {
			data[id] = buildConfigSet(val)
		} else {
			data[id] = d
		}
	}
	return data
}
