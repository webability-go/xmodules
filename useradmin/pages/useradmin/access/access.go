package main

import (
	//	"fmt"
	"strings"

	"encoding/json"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	userassets "github.com/webability-go/xmodules/user/assets"
	"github.com/webability-go/xmodules/user/security"
	"github.com/webability-go/xmodules/useradmin/assets"
)

const (
	ACCESS_SHOWSYSTEM = "ACCESS_SHOWSYSTEM"
)

func Run(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	showsystem := userentries.GetUserParam(ds, userkey, ACCESS_SHOWSYSTEM)
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

func Accessdata(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userkey, _ := ctx.Sessionparams.GetInt("userkey")

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

	showsystem := userentries.GetUserParam(ds, userkey, ACCESS_SHOWSYSTEM)

	filter := xdominion.XConditions{xdominion.NewXCondition("key", "!=", "")}
	if showsystem != "true" {
		filter = append(filter, xdominion.NewXCondition("key", "not like", "\\_%", "and"))
	}
	order := xdominion.XOrder{xdominion.NewXOrderBy("group", xdominion.ASC), xdominion.NewXOrderBy("key", xdominion.ASC)}

	total := userentries.GetAccessesCount(ds, &filter)

	result := map[string]interface{}{
		"total":    total,
		"row":      map[int]xdominion.XRecordDef{},
		"subtotal": 0.1,
		"time":     0.1,
		"subtime":  0.1,
	}

	for _, rg := range dataarray {
		row := rg[0]

		accesses := userentries.GetAccessesList(ds, &filter, &order, rg[1]-rg[0]+1, rg[0])
		if accesses != nil {
			for _, access := range *accesses {
				key, _ := access.GetString("key")
				name, _ := access.GetString("name")
				group, _ := access.GetString("group")
				profiles := searchProfiles(ctx, ds, key)
				users := searchUsers(ctx, ds, key)

				update := ""
				delete := ""
				if key[0] != '_' {
					update = "<button class=\"button update\" onclick=\"useradmin_access_go('" + key + "', 2);\">" + language.Get("button.edit") + "</button>"
					delete = "<button class=\"button delete\" onclick=\"useradmin_access_go('" + key + "', 3);\">" + language.Get("button.delete") + "</button>"
				}

				rec := &xdominion.XRecord{
					"group":    group,
					"key":      key,
					"name":     name,
					"profiles": profiles,
					"users":    users,
					"commands": "<button class=\"button view\" onclick=\"useradmin_access_go('" + key + "', 4);\">" + language.Get("button.view") + "</button> " + update + " " + delete,
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
	userentries.SetUserParam(ds, userkey, ACCESS_SHOWSYSTEM, checked)

	return "OK"
}

func searchProfiles(ctx *context.Context, ds applications.Datasource, key string) string {

	userentries := userassets.GetEntries(ctx.LoggerError)
	recs, err := userentries.GetAccessProfiles(ds, key, 10)
	if recs == nil || err != nil || len(*recs) == 0 {
		return "--"
	}
	result := []string{}
	for _, r := range *recs {
		profilekey, _ := r.GetString("profile")
		result = append(result, "[<b>"+profilekey+"</b>]")
	}
	return strings.Join(result, " ")
}

func searchUsers(ctx *context.Context, ds applications.Datasource, key string) string {

	userentries := userassets.GetEntries(ctx.LoggerError)
	recs, err := userentries.GetAccessUsers(ds, key, 10)
	if recs == nil || err != nil || len(*recs) == 0 {
		return "--"
	}
	result := []string{}
	for _, r := range *recs {
		userkey, _ := r.GetInt("user")
		userdata := userentries.GetUserByKey(ds, userkey)
		username, _ := userdata.GetString("name")
		status, _ := userdata.GetString("status")
		denied, _ := r.GetInt("denied")
		if status == "S" {
			result = append(result, "[<b style=\"color: green;\">"+username+"</b>]")
		} else {
			if denied == 1 {
				result = append(result, "[<b style=\"color: red;\">"+username+"</b>]")
			} else {
				result = append(result, "[<b>"+username+"</b>]")
			}
		}
	}
	return strings.Join(result, " ")
}
