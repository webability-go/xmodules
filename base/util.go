package base

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore/v2"

	"github.com/webability-go/xamboo/assets"
)

func buildTables(ds assets.Datasource) {

	table := baseModule()
	table.SetBase(ds.GetDatabase())
	table.SetLanguage(language.English)
	ds.SetTable("base_module", table)
}

func buildCache(ds assets.Datasource) {

	// Loads all data in XCache
	modules, _ := ds.GetTable("base_module").SelectAll()

	for _, lang := range ds.GetLanguages() {
		canonical := lang.String()
		ds.SetCache("base:modules:"+canonical, xcore.NewXCache("base:modules:"+canonical, 0, 0))

		all := []string{}
		if modules != nil {
			for _, m := range *modules {
				// creates structure on language
				str := CreateStructureModuleByData(ds, m.Clone(), lang)
				key, _ := m.GetString("key")
				all = append(all, key)
				ds.GetCache("base:modules:"+canonical).Set(key, str)
			}
		}
		ds.GetCache("base:modules:"+canonical).Set("all", all)
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
