package material

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"

	"xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["kl_material"] = kl_material()
	sitecontext.Tables["kl_material"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["kl_material"].SetLanguage(language.Spanish)
}

func buildCache(sitecontext *context.Context) {

	// Loads all data in XCache
	materiales, _ := sitecontext.Tables["kl_material"].SelectAll()

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["materiales:"+canonical] = xcore.NewXCache("materiales:"+canonical, 0, 0)

		all := []int{}
		for _, m := range *materiales {
			// creates structure on language
			str := CreateStructureMaterialByData(sitecontext, m.Clone(), lang)
			key, _ := m.GetInt("clave")
			all = append(all, key)
			sitecontext.Caches["materiales:"+canonical].Set(strconv.Itoa(key), str)
		}
		sitecontext.Caches["materiales:"+canonical].Set("all", all)
	}
}

func SynchronizeDatabase(ctx *context.Context) {

	num, err := ctx.Tables["kl_material"].Count(nil)
	if err != nil || num == 0 {
		ctx.Logs["main"].Println("The table kl_material was created (again)")
		ctx.Tables["kl_material"].Synchronize()
	} else {
		ctx.Logs["main"].Println("The table kl_material was not created because it contains data")
	}
}
