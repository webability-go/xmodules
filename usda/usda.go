package usda

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID         = "usda"
	VERSION          = "1.0.0"
	TRANSLATIONTHEME = "nutrient"
)

func InitModule(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	buildCache(sitecontext)
	sitecontext.Modules[MODULEID] = VERSION

	return nil
}

func SynchronizeModule(sitecontext *context.Context, filespath string) []string {

	translation.AddTheme(sitecontext, TRANSLATIONTHEME, "USDA nutrients", translation.SOURCETABLE, "", "name,tag")

	messages := []string{}
	messages = append(messages, createTables(sitecontext)...)

	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := context.AddModule(sitecontext, MODULEID, "List of USDA food and nutrients", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	messages = append(messages, loadTables(sitecontext, filespath)...)
	messages = append(messages, buildCache(sitecontext)...)

	return messages
}

func GetNutrients(sitecontext *context.Context, lang language.Tag) []string {

	canonical := lang.String()

	data, _ := sitecontext.Caches["usda:nutrients:"+canonical].Get("all")
	if data == nil {
		return []string{}
	}
	return data.([]string)
}

func GetNutrient(sitecontext *context.Context, key string, lang language.Tag) *StructureNutrient {

	canonical := lang.String()

	data, _ := sitecontext.Caches["usda:nutrients:"+canonical].Get(key)
	if data == nil {
		sm := CreateStructureNutrientByKey(sitecontext, key, lang)
		if sm == nil {
			sitecontext.Logs["graph"].Println("xmodules::usda::GetNutrient: No hay nutriente creado:", key, lang)
			return nil
		}
		sitecontext.Caches["usda:nutrients:"+canonical].Set(key, sm)
		return sm.(*StructureNutrient)
	}
	return data.(*StructureNutrient)
}

func GetFoodNutrients(sitecontext *context.Context, key string) *xdominion.XRecords {

	nutrients, _ := sitecontext.Tables["usda_foodnutrient"].SelectAll(xdominion.XConditions{
		xdominion.NewXCondition("food", "=", key),
	}, xdominion.NewXOrderBy("nutrient", xdominion.ASC))
	return nutrients
}
