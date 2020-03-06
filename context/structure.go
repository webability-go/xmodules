package context

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	//	"github.com/webability-go/xmodules/translation"
)

// TRANSLATIONTHEME contains the id of the theme to translate the countries
// const TRANSLATIONTHEME = "country"

type StructureModule struct {
	Key  string
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureModuleByKey(sitecontext *Context, key string, lang language.Tag) Structure {
	context_module := sitecontext.GetTable("context_module")
	if context_module == nil {
		sitecontext.Log("main", "Error: the context_module table is not available within the context xmodule")
		return nil
	}
	data, _ := context_module.SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureModuleByData(sitecontext, data, lang)
}

func CreateStructureModuleByData(sitecontext *Context, data xdominion.XRecordDef, lang language.Tag) Structure {

	key, _ := data.GetString("key")

	// builds main data: translations
	/* NOTE: we cannot directly use translation since it's an installed module through this one
	if sitecontext.Tables["country_country"].Language != lang {
		// Only 1 fields to translate: name
		translation.Translate(sitecontext, TRANSLATION_THEME, key, data, map[string]interface{}{"name": true}, sitecontext.Tables["country_country"].Language, lang)
	}
	*/

	return &StructureModule{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sc *StructureModule) ComplementData(sitecontext *Context) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sc *StructureModule) IsAuthorized(sitecontext *Context, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sc *StructureModule) GetData() *xdominion.XRecord {
	return sc.Data
}

// Clone the whole structure
func (sc *StructureModule) Clone() Structure {
	cloned := &StructureModule{
		Key:  sc.Key,
		Lang: sc.Lang,
		Data: sc.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
