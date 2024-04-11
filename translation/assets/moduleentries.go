package assets

import (
	"log"
	"time"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/base"
)

/* Version control:
v0.0.1 a v0.0.4: developement
v0.0.5 - 2023-03-14:
- pages/videoadmin/channel/editor spanish language corrected (channel name description)

*/

const (
	MODULEID       = "kiwitranslation"
	VERSION        = "0.0.3"
	DATASOURCE     = "kiwitranslationdatasource"
	ACCESSGROUP    = "_translationadmin"
	ACCESSLANGUAGE = "_translationadmin_language"
	ACCESSSOURCE   = "_translationadmin_source"
	ACCESSTHEME    = "_translationadmin_theme"
	ACCESS         = "_translationadmin"
	ACCESSTOOLS    = "_translationadmin_tools"
)

var Needs = []string{"base", "user"}

type ModuleEntries struct {
	GetLanguageByKey       func(ds applications.Datasource, key string) *xdominion.XRecord
	GetLanguagesCount      func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetLanguagesList       func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	DeleteLanguageChildren func(ds applications.Datasource, skey string) error
	PruneLanguageChildren  func(ds applications.Datasource, skey string, channel string) error

	GetSourceByKey       func(ds applications.Datasource, key int) *xdominion.XRecord
	GetSourcesCount      func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetSourcesList       func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	DeleteSourceChildren func(ds applications.Datasource, skey string) error
	PruneSourceChildren  func(ds applications.Datasource, skey string, channel string) error

	GetThemeByKey       func(ds applications.Datasource, key int) *xdominion.XRecord
	GetThemeByName      func(ds applications.Datasource, name string) *xdominion.XRecord
	AddTheme            func(ds applications.Datasource, data *xdominion.XRecord) error
	DelThemeByKey       func(ds applications.Datasource, key int) error
	GetThemesCount      func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetThemesList       func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	DeleteThemeChildren func(ds applications.Datasource, skey string) error
	PruneThemeChildren  func(ds applications.Datasource, skey string, channel string) error

	GetTranslationsCount      func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetTranslationsList       func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	DeleteTranslationChildren func(ds applications.Datasource, skey string) error
	PruneTranslationChildren  func(ds applications.Datasource, skey string, channel string) error

	GetTraduccion func(ds applications.Datasource, textooriginal string, tema int, clave string, campo string, lang language.Tag) (string, bool, time.Time, int)
	SetTraduccion func(ds applications.Datasource, textotraducido string, tema int, clave string, campo string, lang language.Tag, verified int) error
}

func GetEntries(logger *log.Logger) *ModuleEntries {
	me := base.GetEntries(logger, MODULEID)
	if me == nil {
		return nil
	}
	lme, ok := me.(*ModuleEntries)
	if !ok {
		return nil
	}
	return lme
}
