package metric

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation"
)

type StructureMetric struct {
	Key  int
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureMetricByKey(sitecontext *base.Datasource, key int, lang language.Tag) base.Structure {
	data, _ := sitecontext.GetTable("metric_unit").SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureMetricByData(sitecontext, data, lang)
}

func CreateStructureMetricByData(sitecontext *base.Datasource, data xdominion.XRecordDef, lang language.Tag) base.Structure {

	key, _ := data.GetInt("key")

	// builds main data: translations
	if sitecontext.GetTable("metric_unit").Language != lang {
		// Only 2 fields to translate: nombre, plural
		translation.Translate(sitecontext, TRANSLATIONTHEME, strconv.Itoa(key), data, map[string]interface{}{"name": true, "plural": true}, sitecontext.GetTable("metric_unit").Language, lang)
	}

	return &StructureMetric{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureMetric) ComplementData(sitecontext *base.Datasource) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureMetric) IsAuthorized(sitecontext *base.Datasource, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructureMetric) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructureMetric) Clone() base.Structure {
	cloned := &StructureMetric{
		Key:  sm.Key,
		Lang: sm.Lang,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
