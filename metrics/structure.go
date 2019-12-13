
package metrics

import (
  "strconv"

  "golang.org/x/text/language"

  "github.com/webability-go/xdominion"

  "xmodules/context"
  "xmodules/translation"
  "xmodules/structure"
)

const(
  TRANSLATION_THEME = 12

)

type StructureMetric struct {
  Key    int
	Lang   language.Tag
	Data   *xdominion.XRecord
}

func CreateStructureMetricByKey(sitecontext *context.Context, key int, lang language.Tag) structure.Structure {
	data, _ := sitecontext.Tables["kl_medida"].SelectOne(key)
	if data == nil {
	  return nil
	}
	return CreateStructureMetricByData(sitecontext, data, lang)
}

func CreateStructureMetricByData(sitecontext *context.Context, data xdominion.XRecordDef, lang language.Tag) structure.Structure {

	key, _ := data.GetInt("clave")

	// builds main data: translations
	if sitecontext.Tables["kl_medida"].Language != lang {
		// Only 2 fields to translate: nombre, plural
		translation.Translate(sitecontext, TRANSLATION_THEME, strconv.Itoa(key), data, map[string]interface{}{"nombre": true, "plural": true}, sitecontext.Tables["kl_medida"].Language, lang)
	}

  return &StructureMetric{Key: key, Lang: lang, Data: data.(*xdominion.XRecord)}
}

	// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureMetric)ComplementData(sitecontext *context.Context) {

}

  // IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureMetric)IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool {
  return true
}

  // Returns the raw data
func (sm *StructureMetric)GetData() *xdominion.XRecord {
  return sm.Data
}

	// Clone the whole structure
func (sm *StructureMetric)Clone() structure.Structure {
  cloned := &StructureMetric{
    Key: sm.Key,
    Lang: sm.Lang,
    Data: sm.Data.Clone().(*xdominion.XRecord),
  }
  return cloned
}
