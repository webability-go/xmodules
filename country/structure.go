package country

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

type StructureCountry struct {
	Key  string
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureCountryByKey(sitecontext *context.Context, key string, lang language.Tag) context.Structure {
	data, _ := sitecontext.GetTable("country_country").SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureCountryByData(sitecontext, data, lang)
}

func CreateStructureCountryByData(sitecontext *context.Context, data xdominion.XRecordDef, lang language.Tag) context.Structure {

	key, _ := data.GetString("key")

	// builds main data: translations
	if sitecontext.GetTable("country_country").Language != lang {
		// Only 1 fields to translate: name
		translation.Translate(sitecontext, TRANSLATIONTHEME, key, data, map[string]interface{}{"name": true}, sitecontext.GetTable("country_country").Language, lang)
	}

	return &StructureCountry{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sc *StructureCountry) ComplementData(sitecontext *context.Context) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sc *StructureCountry) IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sc *StructureCountry) GetData() *xdominion.XRecord {
	return sc.Data
}

// Clone the whole structure
func (sc *StructureCountry) Clone() context.Structure {
	cloned := &StructureCountry{
		Key:  sc.Key,
		Lang: sc.Lang,
		Data: sc.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
