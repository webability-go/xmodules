package main

import (
	"encoding/json"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/cms/context"

	"xmodules/translation/assets"

	"github.com/webability-go/xmodules/base"
	userassets "github.com/webability-go/xmodules/user/assets"
	"github.com/webability-go/xmodules/user/security"
)

const (
	TRANSLATIONLANGUAGE_FILTER = "TRANSLATIONLANGUAGE_FILTER"
)

func Run(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESSLANGUAGE)
	if !ok {
		return ""
	}

	entries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	filter := entries.GetUserParam(ds, userkey, TRANSLATIONLANGUAGE_FILTER)

	params := &xcore.XDataset{
		"SEARCHFILTER": filter,
		"#":            language,
	}

	return template.Execute(params)
}

func Languagedata(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESSLANGUAGE)
	if !ok {
		return ""
	}

	//	userkey, _ := ctx.Sessionparams.GetInt("userkey")
	xdominion.DEBUG = false
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

	search := userentries.GetUserParam(ds, userkey, TRANSLATIONLANGUAGE_FILTER)

	filter := xdominion.XConditions{xdominion.NewXCondition("key", "!=", "")}
	if search != "" {
		filter = append(filter, xdominion.NewXCondition("key", "ilike", "%"+search+"%", "and", 1, 0))
		filter = append(filter, xdominion.NewXCondition("name", "ilike", "%"+search+"%", "or", 0, 1))
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
		order = xdominion.XOrder{xdominion.NewXOrderBy("key", qsort)}
	case "name":
		order = xdominion.XOrder{xdominion.NewXOrderBy("name", qsort)}
	default:
		order = xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	}

	total := entries.GetLanguagesCount(ds, &filter)

	result := map[string]interface{}{
		"total":    total,
		"row":      map[int]xdominion.XRecordDef{},
		"subtotal": 0.1,
		"time":     0.1,
		"subtime":  0.1,
	}

	for _, rg := range dataarray {
		row := rg[0]

		languages := entries.GetLanguagesList(ds, &filter, &order, rg[1]-rg[0]+1, rg[0])
		if languages != nil {
			for _, reclanguage := range *languages {
				key, _ := reclanguage.GetString("key")
				name, _ := reclanguage.GetString("name")
				status, _ := reclanguage.GetString("config")

				update := ""
				delete := ""
				if key[0] != '_' {
					update = "<button class=\"button update\" onclick=\"translationadmin_language_go('" + key + "', 2);\">" + language.Get("button.edit") + "</button>"
					delete = "<button class=\"button delete\" onclick=\"translationadmin_language_go('" + key + "', 3);\">" + language.Get("button.delete") + "</button>"
				}
				// quantity of translations into the language
				quantity := entries.GetTranslationsCount(ds, &xdominion.XConditions{xdominion.NewXCondition("canal", "=", key)})
				// status red/green
				sstatus := "<div style=\"width: 40px; height: 20px; background-color: #aaffaa; padding-top: 6px; text-align: center;\">On</div>"
				if status != "1" {
					sstatus = "<div style=\"width: 40px; height: 20px; background-color: #ffaaaa; padding-top: 6px; text-align: center;\">Off</div>"
				}

				rec := &xdominion.XRecord{
					"clave":    key,
					"nombre":   name,
					"estatus":  sstatus,
					"quantity": quantity,
					"commands": "<button class=\"button view\" onclick=\"translationadmin_language_go('" + key + "', 4);\">" + language.Get("button.view") + "</button> " + update + " " + delete,
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
	userentries.SetUserParam(ds, userkey, TRANSLATIONLANGUAGE_FILTER, filter)

	return "OK"
}
