package ingredient

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

type StructurePasillo struct {
	Key  int
	Lang language.Tag
	Data *xdominion.XRecord
}

type StructureIngredient struct {
	Key  int
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructurePasilloByKey(sitecontext *context.Context, key int, lang language.Tag) context.Structure {
	data, _ := sitecontext.GetTable("kl_ingredientepasillo").SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructurePasilloByData(sitecontext, data, lang)
}

func CreateStructurePasilloByData(sitecontext *context.Context, data xdominion.XRecordDef, lang language.Tag) context.Structure {

	key, _ := data.GetInt("clave")

	// builds main data: translations
	if sitecontext.GetTable("kl_ingredientepasillo").Language != lang {
		// Only 1 fields to translate: nombre
		translation.Translate(sitecontext, TRANSLATIONTHEME, strconv.Itoa(key), data, map[string]interface{}{"nombre": true}, sitecontext.GetTable("kl_ingredientepasillo").Language, lang)
	}

	return &StructurePasillo{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructurePasillo) ComplementData(sitecontext *context.Context) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructurePasillo) IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructurePasillo) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructurePasillo) Clone() context.Structure {
	cloned := &StructurePasillo{
		Key:  sm.Key,
		Lang: sm.Lang,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}

func CreateStructureIngredientByKey(sitecontext *context.Context, key int, lang language.Tag) context.Structure {
	data, _ := sitecontext.GetTable("kl_ingrediente").SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureIngredientByData(sitecontext, data, lang)
}

func CreateStructureIngredientByData(sitecontext *context.Context, data xdominion.XRecordDef, lang language.Tag) context.Structure {

	key, _ := data.GetInt("clave")

	// builds main data: translations
	if sitecontext.GetTable("kl_ingrediente").Language != lang {
		// Only 2 fields to translate: nombre, plural
		translation.Translate(sitecontext, TRANSLATIONTHEMEAISLE, strconv.Itoa(key), data, map[string]interface{}{"nombre": true, "plural": true}, sitecontext.GetTable("kl_ingrediente").Language, lang)
	}

	return &StructureIngredient{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureIngredient) ComplementData(sitecontext *context.Context) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureIngredient) IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructureIngredient) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructureIngredient) Clone() context.Structure {
	cloned := &StructureIngredient{
		Key:  sm.Key,
		Lang: sm.Lang,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
