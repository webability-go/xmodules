package usda

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"

	"xmodules/translation"
)

const (
	TRANSLATION_THEME = 18
)

type StructureNutrient struct {
	Key  string
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureNutrientByKey(sitecontext *context.Context, key string, lang language.Tag) context.Structure {
	data, _ := sitecontext.Tables["usda_nutrient"].SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureNutrientByData(sitecontext, data, lang)
}

func CreateStructureNutrientByData(sitecontext *context.Context, data xdominion.XRecordDef, lang language.Tag) context.Structure {

	key, _ := data.GetString("key")

	// builds main data: translations
	if sitecontext.Tables["usda_nutrient"].Language != lang {
		// Only 2 fields to translate: name, tag
		translation.Translate(sitecontext, TRANSLATION_THEME, key, data, map[string]interface{}{"name": true, "tag": true}, sitecontext.Tables["usda_nutrient"].Language, lang)
	}

	return &StructureNutrient{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureNutrient) ComplementData(sitecontext *context.Context) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureNutrient) IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructureNutrient) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructureNutrient) Clone() context.Structure {
	cloned := &StructureNutrient{
		Key:  sm.Key,
		Lang: sm.Lang,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
