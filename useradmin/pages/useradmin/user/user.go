package main

import (
	"encoding/json"
	"strings"

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
	USER_SHOWALL = "USER_SHOWALL"
)

func Run(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	showall := userentries.GetUserParam(ds, userkey, USER_SHOWALL)
	checked := ""
	if showall == "true" {
		checked = "checked "
	}

	params := &xcore.XDataset{
		"CHECKED": checked,
		"#":       language,
	}

	return template.Execute(params)
}

func Userdata(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

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

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	showall := userentries.GetUserParam(ds, userkey, USER_SHOWALL)

	filter := xdominion.XConditions{xdominion.NewXCondition("key", "!=", 0)}
	if showall != "true" {
		filter = append(filter, xdominion.NewXCondition("status", "in", "('A', 'S')", "and"))
	}
	order := xdominion.XOrder{xdominion.NewXOrderBy("name", xdominion.ASC)}

	total := userentries.GetUsersCount(ds, &filter)

	result := map[string]interface{}{
		"total":    total,
		"row":      map[int]xdominion.XRecordDef{},
		"subtotal": 0.1,
		"time":     0.1,
		"subtime":  0.1,
	}

	for _, rg := range dataarray {
		row := rg[0]

		users := userentries.GetUsersList(ds, &filter, &order, rg[1]-rg[0]+1, rg[0])
		if users != nil {
			for _, user := range *users {
				key, _ := user.GetInt("key")
				skey, _ := user.GetString("key")
				name, _ := user.GetString("name")
				status, _ := user.GetString("status")
				accesses := searchAccesses(ctx, ds, key)
				profiles := searchProfiles(ctx, ds, key)

				update := "<button class=\"button update\" onclick=\"useradmin_user_go(" + skey + ", 2);\">" + language.Get("button.edit") + "</button>"
				delete := "<button class=\"button delete\" onclick=\"useradmin_user_go(" + skey + ", 3);\">" + language.Get("button.delete") + "</button>"

				rec := &xdominion.XRecord{
					"key":      key,
					"status":   status,
					"name":     name,
					"accesses": accesses,
					"profiles": profiles,
					"commands": "<button class=\"button view\" onclick=\"useradmin_user_go(" + skey + ", 4);\">" + language.Get("button.view") + "</button> " + update + " " + delete,
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

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	checked := ctx.Request.Form.Get("checked")
	userentries.SetUserParam(ds, userkey, USER_SHOWALL, checked)

	return "OK"
}

func searchAccesses(ctx *context.Context, ds applications.Datasource, key int) string {

	userentries := userassets.GetEntries(ctx)
	recs, err := userentries.GetUserAccesses(ds, key, 10)
	if recs == nil || err != nil || len(*recs) == 0 {
		return "--"
	}
	result := []string{}
	for _, r := range *recs {
		accesskey, _ := r.GetString("access")
		denied, _ := r.GetInt("denied")
		if denied == 1 {
			result = append(result, "[<b style=\"color: red;\">"+accesskey+"</b>]")
		} else {
			result = append(result, "[<b>"+accesskey+"</b>]")
		}
	}
	return strings.Join(result, "")
}

func searchProfiles(ctx *context.Context, ds applications.Datasource, key int) string {

	userentries := userassets.GetEntries(ctx)
	recs, err := userentries.GetUserProfiles(ds, key, 10)
	if recs == nil || err != nil || len(*recs) == 0 {
		return "--"
	}
	result := []string{}
	for _, r := range *recs {
		profilekey, _ := r.GetString("profile")
		result = append(result, "[<b>"+profilekey+"</b>]")
	}
	return strings.Join(result, "")
}
