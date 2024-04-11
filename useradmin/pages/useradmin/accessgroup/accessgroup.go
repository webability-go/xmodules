package main

import (
	"encoding/json"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	userassets "github.com/webability-go/xmodules/user/assets"
	"github.com/webability-go/xmodules/user/security"
	"github.com/webability-go/xmodules/useradmin/assets"
)

const (
	ACCESSGROUP_SHOWSYSTEM = "ACCESSGROUP_SHOWSYSTEM"
)

func Run(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	showsystem := userentries.GetUserParam(ds, userkey, ACCESSGROUP_SHOWSYSTEM)
	checked := ""
	if showsystem == "true" {
		checked = "checked "
	}

	params := &xcore.XDataset{
		"CHECKED": checked,
		"#":       language,
	}

	return template.Execute(params)
}

func Accessgroupdata(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	var dataarray [][]int
	data := ctx.Request.Form.Get("data")
	if data != "" {
		json.Unmarshal([]byte(data), &dataarray)
	}
	if dataarray == nil {
		dataarray = [][]int{{0, 49}}
	}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	showsystem := userentries.GetUserParam(ds, userkey, ACCESSGROUP_SHOWSYSTEM)

	filter := xdominion.XConditions{xdominion.NewXCondition("key", "!=", "")}
	if showsystem != "true" {
		filter = append(filter, xdominion.NewXCondition("key", "not like", "\\_%", "and"))
	}
	order := xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}

	total := userentries.GetAccessGroupsCount(ds, &filter)

	result := map[string]interface{}{
		"total": total,
		"row":   map[int]xdominion.XRecordDef{},
		"time":  0.1,
	}

	for _, rg := range dataarray {
		row := rg[0]

		accessgroups := userentries.GetAccessGroupsList(ds, &filter, &order, rg[1]-rg[0]+1, rg[0])
		if accessgroups != nil {
			for _, accessgroup := range *accessgroups {
				key, _ := accessgroup.GetString("key")
				name, _ := accessgroup.GetString("name")

				update := ""
				delete := ""
				if key[0] != '_' {
					update = "<button class=\"button update\" onclick=\"useradmin_accessgroup_go('" + key + "', 2);\">" + language.Get("button.edit") + "</button>"
					delete = "<button class=\"button delete\" onclick=\"useradmin_accessgroup_go('" + key + "', 3);\">" + language.Get("button.delete") + "</button>"
				}

				rec := &xdominion.XRecord{
					"key":      key,
					"name":     name,
					"commands": "<button class=\"button view\" onclick=\"useradmin_accessgroup_go('" + key + "', 4);\">" + language.Get("button.view") + "</button> " + update + " " + delete,
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
	userentries.SetUserParam(ds, userkey, ACCESSGROUP_SHOWSYSTEM, checked)

	return "OK"
}
