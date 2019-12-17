package country

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

// TRANSLATIONTHEME contains the id of the theme to translate the countries
const TRANSLATIONTHEME = "country"

type StructureCountry struct {
	Key  string
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureCountryByKey(sitecontext *context.Context, key string, lang language.Tag) structure.Structure {
	data, _ := sitecontext.Tables["country_country"].SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureCountryByData(sitecontext, data, lang)
}

func CreateStructureCountryByData(sitecontext *context.Context, data xdominion.XRecordDef, lang language.Tag) structure.Structure {

	key, _ := data.GetString("key")

	// builds main data: translations
	if sitecontext.Tables["country_country"].Language != lang {
		// Only 1 fields to translate: name
		translation.Translate(sitecontext, TRANSLATION_THEME, key, data, map[string]interface{}{"name": true}, sitecontext.Tables["country_country"].Language, lang)
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
func (sc *StructureCountry) Clone() structure.Structure {
	cloned := &StructureCountry{
		Key:  sc.Key,
		Lang: sc.Lang,
		Data: sc.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
