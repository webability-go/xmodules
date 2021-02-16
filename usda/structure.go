package usda

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation"

	"github.com/webability-go/xamboo/applications"
)

type StructureNutrient struct {
	Key  string
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureNutrientByKey(sitecontext applications.Datasource, key string, lang language.Tag) base.Structure {
	data, _ := sitecontext.GetTable("usda_nutrient").SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureNutrientByData(sitecontext, data, lang)
}

func CreateStructureNutrientByData(sitecontext applications.Datasource, data xdominion.XRecordDef, lang language.Tag) base.Structure {

	key, _ := data.GetString("key")

	// builds main data: translations
	if sitecontext.GetTable("usda_nutrient").Language != lang {
		// Only 2 fields to translate: name, tag
		translation.Translate(sitecontext, TRANSLATIONTHEME, key, data, map[string]interface{}{"name": true, "tag": true}, sitecontext.GetTable("usda_nutrient").Language, lang)
	}

	return &StructureNutrient{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureNutrient) ComplementData(sitecontext applications.Datasource) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureNutrient) IsAuthorized(sitecontext applications.Datasource, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructureNutrient) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructureNutrient) Clone() base.Structure {
	cloned := &StructureNutrient{
		Key:  sm.Key,
		Lang: sm.Lang,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
