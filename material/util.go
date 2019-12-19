package material

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["material_material"] = materialMaterial()
	sitecontext.Tables["material_material"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["material_material"].SetLanguage(language.Spanish)
}

func buildCache(sitecontext *context.Context) {

	// Loads all data in XCache
	materiales, _ := sitecontext.Tables["material_material"].SelectAll()
	if materiales == nil {
		return
	}

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
