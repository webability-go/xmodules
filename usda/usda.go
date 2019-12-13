package usda

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"

	"xmodules/context"
)

func InitUSDA(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	buildCache(sitecontext)

	return nil
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
