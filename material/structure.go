package material

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

type StructureMaterial struct {
	Key  int
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureMaterialByKey(sitecontext *context.Context, key int, lang language.Tag) context.Structure {
	data, _ := sitecontext.Tables["kl_material"].SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureMaterialByData(sitecontext, data, lang)
}

func CreateStructureMaterialByData(sitecontext *context.Context, data xdominion.XRecordDef, lang language.Tag) context.Structure {

	key, _ := data.GetInt("clave")

	// builds main data: translations
	if sitecontext.Tables["kl_material"].Language != lang {
		// Only 1 fields to translate: nombre
		translation.Translate(sitecontext, TRANSLATIONTHEME, strconv.Itoa(key), data, map[string]interface{}{"nombre": true, "plural": true}, sitecontext.Tables["kl_material"].Language, lang)
	}

	return &StructureMaterial{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureMaterial) ComplementData(sitecontext *context.Context) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureMaterial) IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructureMaterial) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructureMaterial) Clone() context.Structure {
	cloned := &StructureMaterial{
		Key:  sm.Key,
		Lang: sm.Lang,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
