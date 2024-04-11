package main

import (
	"encoding/json"
	"fmt"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/cms/context"

	"xmodules/translation/assets"

	"github.com/webability-go/xmodules/base"
	userassets "github.com/webability-go/xmodules/user/assets"
	"github.com/webability-go/xmodules/user/security"
)

const (
	TRANSLATIONTHEME_FILTER = "TRANSLATIONTHEME_FILTER"
)

func Run(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESSTHEME)
	if !ok {
		return ""
	}

	entries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	filter := entries.GetUserParam(ds, userkey, TRANSLATIONTHEME_FILTER)

	params := &xcore.XDataset{
		"SEARCHFILTER": filter,
		"#":            language,
	}

	return template.Execute(params)
}

func Themedata(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESSTHEME)
	if !ok {
		return ""
	}

	//	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	var dataarray [][]int
	data := ctx.Request.Form.Get("data")
	field := ctx.Request.Form.Get("field")
	sort := ctx.Request.Form.Get("sort")

	if data != "" {
		json.Unmarshal([]byte(data), &dataarray)
	}
	if dataarray == nil {
		dataarray = [][]int{{0, 49}}
	}

	entries := assets.GetEntries(ctx.LoggerError)
	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	search := userentries.GetUserParam(ds, userkey, TRANSLATIONTHEME_FILTER)

	filter := xdominion.XConditions{xdominion.NewXCondition("clave", "!=", 0)}
	if search != "" {
		filter = append(filter, xdominion.NewXCondition("clave", "ilike", "%"+search+"%", "and", 1, 0))
		filter = append(filter, xdominion.NewXCondition("nombre", "ilike", "%"+search+"%", "or", 0, 1))
	}
	qsort := xdominion.ASC
	if sort == "desc" {
		qsort = xdominion.DESC
	} else if sort != "asc" {
		field = ""
	}

	var order xdominion.XOrder
	switch field {
	case "key":
		order = xdominion.XOrder{xdominion.NewXOrderBy("clave", qsort)}
	case "name":
		order = xdominion.XOrder{xdominion.NewXOrderBy("nombre", qsort)}
	default:
		order = xdominion.XOrder{xdominion.NewXOrderBy("clave", xdominion.ASC)}
	}

	total := entries.GetThemesCount(ds, &filter)

	result := map[string]interface{}{
		"total":    total,
		"row":      map[int]xdominion.XRecordDef{},
		"subtotal": 0.1,
		"time":     0.1,
		"subtime":  0.1,
	}

	sources := GetSourcesList(language)

	for _, rg := range dataarray {
		row := rg[0]

		themes := entries.GetThemesList(ds, &filter, &order, rg[1]-rg[0]+1, rg[0])
		if themes != nil {
			for _, theme := range *themes {
				key, _ := theme.GetString("clave")
				name, _ := theme.GetString("nombre")
				source, _ := theme.GetInt("fuente")
				config, _ := theme.GetString("config")

				update := ""
				delete := ""
				if key[0] != '_' {
					update = "<button class=\"button update\" onclick=\"translationadmin_theme_go('" + key + "', 2);\">" + language.Get("button.edit") + "</button>"
					delete = "<button class=\"button delete\" onclick=\"translationadmin_theme_go('" + key + "', 3);\">" + language.Get("button.delete") + "</button>"
				}
				// quantity of translations into the theme
				quantity := entries.GetTranslationsCount(ds, &xdominion.XConditions{xdominion.NewXCondition("tema", "=", key)})

				rec := &xdominion.XRecord{
					"clave":    key,
					"nombre":   name,
					"source":   sources[source],
					"config":   language.Get("quantity") + " <b>" + fmt.Sprintf("%d", quantity) + "</b><br />" + config,
					"commands": "<button class=\"button view\" onclick=\"translationadmin_theme_go('" + key + "', 4);\">" + language.Get("button.view") + "</button> " + update + " " + delete,
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

	filter := ctx.Request.Form.Get("filter")
	userentries.SetUserParam(ds, userkey, TRANSLATIONTHEME_FILTER, filter)

	return "OK"
}

func GetSourcesList(language *xcore.XLanguage) map[int]string {
	return map[int]string{
		1:  language.Get("source.1"),
		2:  language.Get("source.2"),
		3:  language.Get("source.3"),
		20: language.Get("source.20"),
	}
}
