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
	PROFILE_SHOWALL = "PROFILE_SHOWALL"
)

func Run(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	showall := userentries.GetUserParam(ds, userkey, PROFILE_SHOWALL)
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

func Profiledata(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

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

	showall := userentries.GetUserParam(ds, userkey, PROFILE_SHOWALL)

	filter := xdominion.XConditions{xdominion.NewXCondition("key", "!=", "")}
	if showall != "true" {
		filter = append(filter, xdominion.NewXCondition("status", "=", 1, "and"))
	}
	order := xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}

	total := userentries.GetProfilesCount(ds, &filter)

	result := map[string]interface{}{
		"total":    total,
		"row":      map[int]xdominion.XRecordDef{},
		"subtotal": 0.1,
		"time":     0.1,
		"subtime":  0.1,
	}

	for _, rg := range dataarray {
		row := rg[0]

		profiles := userentries.GetProfilesList(ds, &filter, &order, rg[1]-rg[0]+1, rg[0])
		if profiles != nil {
			for _, profile := range *profiles {
				key, _ := profile.GetString("key")
				name, _ := profile.GetString("name")
				status, _ := profile.GetString("status")
				sstatus := ""
				if status == "1" {
					sstatus = language.Get("status.yes")
				} else {
					sstatus = language.Get("status.no")
				}
				accesses := searchAccesses(ctx, ds, key)
				users := searchUsers(ctx, ds, key)

				update := ""
				delete := ""
				if key[0] != '_' {
					update = "<button class=\"button update\" onclick=\"useradmin_profile_go('" + key + "', 2);\">" + language.Get("button.edit") + "</button>"
					delete = "<button class=\"button delete\" onclick=\"useradmin_profile_go('" + key + "', 3);\">" + language.Get("button.delete") + "</button>"
				}

				rec := &xdominion.XRecord{
					"key":      key,
					"status":   sstatus,
					"name":     name,
					"accesses": accesses,
					"users":    users,
					"commands": "<button class=\"button view\" onclick=\"useradmin_profile_go('" + key + "', 4);\">" + language.Get("button.view") + "</button> " + update + " " + delete,
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
	userentries.SetUserParam(ds, userkey, PROFILE_SHOWALL, checked)

	return "OK"
}

func searchAccesses(ctx *context.Context, ds applications.Datasource, key string) string {

	userentries := userassets.GetEntries(ctx)
	recs, err := userentries.GetProfileAccesses(ds, key, 10)
	if recs == nil || err != nil || len(*recs) == 0 {
		return "--"
	}
	result := []string{}
	for _, r := range *recs {
		accesskey, _ := r.GetString("access")
		result = append(result, "[<b>"+accesskey+"</b>]")
	}
	return strings.Join(result, " ")
}

func searchUsers(ctx *context.Context, ds applications.Datasource, key string) string {

	userentries := userassets.GetEntries(ctx)
	recs, err := userentries.GetProfileUsers(ds, key, 10)
	if recs == nil || err != nil || len(*recs) == 0 {
		return "--"
	}
	result := []string{}
	for _, r := range *recs {
		userkey, _ := r.GetInt("user")
		userdata := userentries.GetUserByKey(ds, userkey)
		username, _ := userdata.GetString("name")
		result = append(result, "[<b>"+username+"</b>]")
	}
	return strings.Join(result, " ")
}
