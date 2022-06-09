package country

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation"

	"github.com/webability-go/xamboo/applications"
)

type StructureCountry struct {
	Key  string
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureCountryByKey(sitecontext applications.Datasource, key string, lang language.Tag) base.Structure {
	data, _ := sitecontext.GetTable("country_country").SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureCountryByData(sitecontext, data, lang)
}

func CreateStructureCountryByData(sitecontext applications.Datasource, data xdominion.XRecordDef, lang language.Tag) base.Structure {

	key, _ := data.GetString("key")

	// builds main data: translations
	if sitecontext.GetTable("country_country").Language != lang {
		// Only 1 fields to translate: name
		translation.Translate(sitecontext, TRANSLATIONTHEME, key, data, map[string]interface{}{"name": true}, sitecontext.GetTable("country_country").Language, lang)
	}

	return &StructureCountry{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sc *StructureCountry) ComplementData(sitecontext applications.Datasource) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sc *StructureCountry) IsAuthorized(sitecontext applications.Datasource, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sc *StructureCountry) GetData() *xdominion.XRecord {
	return sc.Data
}

// Clone the whole structure
func (sc *StructureCountry) Clone() base.Structure {
	cloned := &StructureCountry{
		Key:  sc.Key,
		Lang: sc.Lang,
		Data: sc.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
