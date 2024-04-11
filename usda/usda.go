package usda

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xdominion"
)

func GetNutrients(ds applications.Datasource, lang language.Tag) []string {

	canonical := lang.String()

	cache := ds.GetCache("usda:nutrients:" + canonical)
	if cache == nil {
		ds.Log("main", "xmodules::usda::GetNutrients: Error, the nutrients cache is not available on this site context")
		return nil
	}

	data, _ := cache.Get("all")
	if data == nil {
		return []string{}
	}
	return data.([]string)
}

func GetNutrient(ds applications.Datasource, key string, lang language.Tag) *StructureNutrient {

	canonical := lang.String()

	cache := ds.GetCache("usda:nutrients:" + canonical)
	if cache == nil {
		ds.Log("main", "xmodules::usda::GetNutrient: Error, the nutrients cache is not available on this site context")
		return nil
	}

	data, _ := cache.Get(key)
	if data == nil {
		sm := CreateStructureNutrientByKey(ds, key, lang)
		if sm == nil {
			ds.Log("graph", "xmodules::usda::GetNutrient: No hay nutriente creado:", key, lang)
			return nil
		}
		cache.Set(key, sm)
		return sm.(*StructureNutrient)
	}
	return data.(*StructureNutrient)
}

func GetFoodNutrients(ds applications.Datasource, key string) *xdominion.XRecords {

	usda_foodnutrient := ds.GetTable("usda_foodnutrient")
	if usda_foodnutrient == nil {
		ds.Log("main", "xmodules::usda::GetFoodNutrients: Error, the usda_foodnutrient table is not available on this context")
		return nil
	}

	nutrients, _ := usda_foodnutrient.SelectAll(xdominion.XConditions{
		xdominion.NewXCondition("food", "=", key),
	}, xdominion.NewXOrderBy("nutrient", xdominion.ASC))
	return nutrients
}
