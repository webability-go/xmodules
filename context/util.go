package context

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore"
)

func buildTables(sitecontext *Context, databasename string) {

	table := contextModule()
	table.SetBase(sitecontext.GetDatabase(databasename))
	table.SetLanguage(language.English)
	sitecontext.SetTable("context_module", table)
}

func buildCache(sitecontext *Context) {

	// Loads all data in XCache
	modules, _ := sitecontext.GetTable("context_module").SelectAll()

	for _, lang := range sitecontext.GetLanguages() {
		canonical := lang.String()
		sitecontext.SetCache("context:modules:"+canonical, xcore.NewXCache("context:modules:"+canonical, 0, 0))

		all := []string{}
		if modules != nil {
			for _, m := range *modules {
				// creates structure on language
				str := CreateStructureModuleByData(sitecontext, m.Clone(), lang)
				key, _ := m.GetString("key")
				all = append(all, key)
				sitecontext.GetCache("context:modules:"+canonical).Set(key, str)
			}
		}
		sitecontext.GetCache("context:modules:"+canonical).Set("all", all)
	}
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
