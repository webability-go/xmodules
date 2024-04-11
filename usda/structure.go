package usda

import (
	"fmt"
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/base"

	//	"github.com/webability-go/xmodules/translation"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/translation"
	"github.com/webability-go/xmodules/usda/assets"
)

var USDA_TRADUCCIONCAMPOS = map[string]interface{}{
	"name": true,
	"tag":  true,
} // campos necesarios a traducir
type StructureNutrient struct {
	Key  string
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureNutrientByKey(ds applications.Datasource, key string, lang language.Tag) base.Structure {
	data, _ := ds.GetTable("usda_nutrient").SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureNutrientByData(ds, data, lang)
}

func CreateStructureNutrientByData(ds applications.Datasource, data xdominion.XRecordDef, lang language.Tag) base.Structure {

	key, _ := data.GetString("key")

	if ds.GetTable("usda_nutrient").Language == lang {
		return &StructureNutrient{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
	}

	tema := 0
	prompt := ""
	config := ds.GetConfig()
	if config != nil {
		localconfig := config.GetConfig(assets.MODULEID)
		if localconfig != nil {
			stema, _ := localconfig.GetString(assets.TRANSLATIONTHEME)
			tema, _ = strconv.Atoi(stema)
			prompt, _ = localconfig.GetString(assets.PROMPT_USDA)
		}
	}
	if prompt == "" || tema == 0 {
		fmt.Println("ERROR: FALTA PROMPT Y/O TEMA EN xmodules/usda/structure ", tema)
	}

	translation.TranslatePrompt(ds, tema, key, data, USDA_TRADUCCIONCAMPOS, ds.GetTable("usda_nutrient").Language, lang, prompt)

	return &StructureNutrient{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureNutrient) ComplementData(ds applications.Datasource) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureNutrient) IsAuthorized(ds applications.Datasource, site string, language string, device string) bool {
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
