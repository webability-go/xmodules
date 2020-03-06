package metric

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

type StructureMetric struct {
	Key  int
	Lang language.Tag
	Data *xdominion.XRecord
}

func CreateStructureMetricByKey(sitecontext *context.Context, key int, lang language.Tag) context.Structure {
	data, _ := sitecontext.GetTable("metric_unit").SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureMetricByData(sitecontext, data, lang)
}

func CreateStructureMetricByData(sitecontext *context.Context, data xdominion.XRecordDef, lang language.Tag) context.Structure {

	key, _ := data.GetInt("key")

	// builds main data: translations
	if sitecontext.GetTable("metric_unit").Language != lang {
		// Only 2 fields to translate: nombre, plural
		translation.Translate(sitecontext, TRANSLATIONTHEME, strconv.Itoa(key), data, map[string]interface{}{"name": true, "plural": true}, sitecontext.GetTable("metric_unit").Language, lang)
	}

	return &StructureMetric{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureMetric) ComplementData(sitecontext *context.Context) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureMetric) IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructureMetric) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructureMetric) Clone() context.Structure {
	cloned := &StructureMetric{
		Key:  sm.Key,
		Lang: sm.Lang,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
