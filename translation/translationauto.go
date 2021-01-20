package translation

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	//	"github.com/webability-go/xmodules/base"

	serverassets "github.com/webability-go/xamboo/assets"
)

// All purpose Translation method for structures
// Any MAINDATA record MUST have a lastmodif field
func Translate(sitecontext serverassets.Datasource, theme string, key string, maindata xdominion.XRecordDef, fields map[string]interface{}, fromLang language.Tag, toLang language.Tag) {

	lastmodif, _ := maindata.GetTime("lastmodif")
	trtbl := NewTranslationBlock(theme, key, lastmodif, fromLang, toLang)

	for campo, sub := range fields {
		val := ""
		switch sub.(type) {
		case bool, int, string:
			val, _ = maindata.GetString(campo)
			if val == "" {
				continue
			}
			trtbl.Set(campo, val)

		case map[string]string:
			subdata, _ := maindata.Get(campo)
			if subdata == nil {
				continue
			}
			switch subdata.(type) {
			case *xdominion.XRecords:
				for _, subrecord := range *subdata.(*xdominion.XRecords) {
					for subcampo, prefix := range sub.(map[string]string) {
						subval, _ := subrecord.GetString(subcampo)
						if subval == "" {
							continue
						}
						subclave, _ := subrecord.GetString("clave")
						if subclave == "" {
							subclave, _ = subrecord.GetString("key")
						}
						trtbl.Set(prefix+subclave, subval)
					}
				}
			}
		}
	}

	trtbl.Verify(sitecontext)

	for campo, sub := range fields {
		switch sub.(type) {
		case bool, int, string:
			maindata.Set(campo, trtbl.Get(campo))
		case map[string]string:
			subdata, _ := maindata.Get(campo)
			if subdata == nil {
				continue
			}
			switch subdata.(type) {
			case *xdominion.XRecords:
				for _, subrecord := range *subdata.(*xdominion.XRecords) {
					for subcampo, prefix := range sub.(map[string]string) {
						subclave, _ := subrecord.GetString("clave")
						if subclave == "" {
							subclave, _ = subrecord.GetString("key")
						}
						subrecord.Set(subcampo, trtbl.Get(prefix+subclave))
					}
				}
			}
		}
	}
}
