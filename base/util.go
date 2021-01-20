package base

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore/v2"

	"github.com/webability-go/xamboo/assets"
)

func linkTables(ds assets.Datasource) {

	table := baseModule()
	table.SetBase(ds.GetDatabase())
	table.SetLanguage(language.English)
	ds.SetTable("base_module", table)
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
