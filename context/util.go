package context

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore/v2"
)

func buildTables(ctx *Context) {

	table := contextModule()
	table.SetBase(ctx.GetDatabase())
	table.SetLanguage(language.English)
	ctx.SetTable("context_module", table)
}

func buildCache(ctx *Context) {

	// Loads all data in XCache
	modules, _ := ctx.GetTable("context_module").SelectAll()

	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()
		ctx.SetCache("context:modules:"+canonical, xcore.NewXCache("context:modules:"+canonical, 0, 0))

		all := []string{}
		if modules != nil {
			for _, m := range *modules {
				// creates structure on language
				str := CreateStructureModuleByData(ctx, m.Clone(), lang)
				key, _ := m.GetString("key")
				all = append(all, key)
				ctx.GetCache("context:modules:"+canonical).Set(key, str)
			}
		}
		ctx.GetCache("context:modules:"+canonical).Set("all", all)
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
