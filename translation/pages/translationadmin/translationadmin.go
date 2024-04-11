package main

import (
	"encoding/json"

	//	"strings"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	//	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	userassets "github.com/webability-go/xmodules/user/assets"
	"github.com/webability-go/xmodules/user/security"

	"xmodules/translation/assets"
)

const (
	TRANSLATION_SHOWALL = "TRANSLATION_SHOWALL"
	TRANSLATION_FILTER  = "TRANSLATION_FILTER"
	TRANSLATION_THEME   = "TRANSLATION_THEME"
)

func Run(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	entries := assets.GetEntries(ctx.LoggerError)
	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	checked := ""
	showall := userentries.GetUserParam(ds, userkey, TRANSLATION_SHOWALL)
	filter := userentries.GetUserParam(ds, userkey, TRANSLATION_FILTER)
	if showall == "true" {
		checked = "checked "
	}
	// Options
	theme := userentries.GetUserParam(ds, userkey, TRANSLATION_THEME)
	order := xdominion.XOrder{xdominion.NewXOrderBy("nombre", xdominion.ASC)}
	options := entries.GetThemesList(ds, nil, &order, 10000000, 0)

	textoptions := ""
	for _, option := range *options {
		optionkey, _ := option.GetString("clave")
		optionname, _ := option.GetString("nombre")
		status, _ := option.GetInt("estatus")
		selected := ""
		if optionkey == theme {
			selected = " selected"
		}
		ostatus := ""
		if status == 1 {
			ostatus = ""
		} else {
			ostatus = " style=\"background-color: #fcc;\""
		}
		textoptions += "<option" + ostatus + selected + " value=\"" + optionkey + "\">" + optionkey + " / " + optionname + "</option>"
	}

	params := &xcore.XDataset{
		"CHECKED":      checked,
		"SEARCHFILTER": filter,
		"THEMEFILTER":  theme,
		"OPTIONS":      textoptions,
		"#":            language,
	}

	return template.Execute(params)
}

func Translationdata(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	var dataarray [][]int
	data := ctx.Request.Form.Get("data")

	// Get all the data of the theme
	// theme := ctx.Request.Form.Get("theme")

	if data != "" {
		json.Unmarshal([]byte(data), &dataarray)
	}
	if dataarray == nil {
		dataarray = [][]int{{0, 49}}
	}

	entries := assets.GetEntries(ctx.LoggerError)
	//	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)
	//	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	//	showall := userentries.GetUserParam(ds, userkey, TRANSLATION_SHOWALL)
	//	search := userentries.GetUserParam(ds, userkey, TRANSLATION_FILTER)
	//	theme := userentries.GetUserParam(ds, userkey, TRANSLATION_THEME)
	theme := 1

	themedata := entries.GetThemeByKey(ds, theme)
	source, _ := themedata.GetInt("fuente")
	name, _ := themedata.GetString("name")

	if source != 20 { // .language file for GO Kiwi7
		result := map[string]interface{}{
			"total":    0,
			"row":      map[int]xdominion.XRecordDef{},
			"subtotal": 0.1,
			"time":     0.1,
			"subtime":  0.1,
		}
		return result
	}

	getfile(name, "en")
	total := 1

	filter := xdominion.XConditions{}
	order := xdominion.XOrder{xdominion.NewXOrderBy("clave", xdominion.ASC)}

	result := map[string]interface{}{
		"total":    total,
		"row":      map[int]xdominion.XRecordDef{},
		"subtotal": 0.1,
		"time":     0.1,
		"subtime":  0.1,
	}

	for _, rg := range dataarray {
		row := rg[0]

		translations := entries.GetTranslationsList(ds, &filter, &order, rg[1]-rg[0]+1, rg[0])
		if translations != nil {
			for _, translation := range *translations {
				key, _ := translation.GetInt("clave")
				skey, _ := translation.GetString("clave")
				theme, _ := translation.GetString("tema")
				lang, _ := translation.GetString("idioma")
				field, _ := translation.GetString("campo")
				translation, _ := translation.GetString("traduccion")

				update := "<button class=\"button update\" onclick=\"translationadmin_translation_go(" + skey + ", 2);\">" + language.Get("button.edit") + "</button>"
				delete := "<button class=\"button delete\" onclick=\"translationadmin_translation_go(" + skey + ", 3);\">" + language.Get("button.delete") + "</button>"

				rec := &xdominion.XRecord{
					"key":         key,
					"theme":       theme,
					"lang":        lang,
					"field":       field,
					"translation": translation,
					"commands":    "<button class=\"button view\" onclick=\"translationadmin_translation_go(" + skey + ", 4);\">" + language.Get("button.view") + "</button> " + update + " " + delete,
				}

				result["row"].(map[int]xdominion.XRecordDef)[row] = rec
				row++
			}
		}
	}

	return result
}

func Filter(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	checked := ctx.Request.Form.Get("checked")
	userentries.SetUserParam(ds, userkey, TRANSLATION_SHOWALL, checked)
	filter := ctx.Request.Form.Get("filter")
	userentries.SetUserParam(ds, userkey, TRANSLATION_FILTER, filter)
	theme := ctx.Request.Form.Get("theme")
	userentries.SetUserParam(ds, userkey, TRANSLATION_THEME, theme)

	return "OK"
}

// read the language
func getfile(path string, language string) {

}
