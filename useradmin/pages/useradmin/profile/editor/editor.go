package main

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xdommask"

	//	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
	userassets "github.com/webability-go/xmodules/user/assets"
	"github.com/webability-go/xmodules/user/security"
	"github.com/webability-go/xmodules/useradmin/assets"
)

var language *xcore.XLanguage

func Run(ctx *context.Context, template *xcore.XTemplate, xlanguage *xcore.XLanguage, e interface{}) interface{} {

	if language == nil {
		language = xlanguage
	}

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	mode := ctx.Request.Form.Get("mode")
	key := ctx.Request.Form.Get("Key")
	if key == "" {
		key = "new"
	}

	params := &xcore.XDataset{
		"FORMPROFILE": createXMLMask("formprofile", mode, ctx),
		"KEY":         key,
		"#":           language,
	}

	return template.Execute(params)
}

func createMask(id string, ctx *context.Context) (*xdommask.Mask, error) {

	hooks := xdommask.MaskHooks{
		Build:     build,
		PreDelete: predelete,
	}
	return xdommask.NewMask(id, hooks, ctx)
}

func createXMLMask(id string, mode string, ctx *context.Context) string {
	mask, _ := createMask(id, ctx)
	cmask := mask.Compile(mode, ctx)
	xmlmask, _ := xml.Marshal(cmask)
	return string(xmlmask)
}

func build(mask *xdommask.Mask, ctx *context.Context) error {

	// Check security
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	if ds == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errordatasource", "build", userassets.DATASOURCE))
	}

	user_profile := ds.GetTable("user_profile")
	if user_profile == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errortable", "build", "user_profile"))
	}

	mask.Table = user_profile
	key := ctx.Request.Form.Get(mask.VarKey)
	if key != "" {
		mask.Key = key
	}

	mask.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	mask.KeyField = "key"

	mask.AlertMessage = "##mask.errormessage##"
	mask.ServerMessage = "##mask.servermessage##"
	mask.InsertTitle = "##mask.titleinsert##"
	mask.DoInsertMessage = "##mask.titleinsert.message##"
	mask.UpdateTitle = "##mask.titleupdate##"
	mask.DoUpdateMessage = "##mask.titleupdate.message##"
	mask.DeleteTitle = "##mask.titledelete##"
	mask.DoDeleteMessage = "##mask.titledelete.message##"
	mask.ViewTitle = "##mask.titleview##"
	mask.FailureJS = "function(params) { this.icontainer.setMessages(params); }"

	// key
	f10 := xdommask.NewTextField("key")
	f10.Title = "##key.title##"
	f10.HelpDescription = "##key.help.description##"
	f10.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f10.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f10.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f10.ViewModes = xdommask.DELETE | xdommask.VIEW
	f10.StatusNotNull = "##key.status.notnull##"
	f10.MaxLength = 30
	f10.Size = "400"
	f10.URLVariable = "key"
	f10.Format = "^[a-z][a-z0-9-_]{1,29}$"
	f10.FormatJS = "^[a-z][a-z0-9-_]{1,29}$"
	f10.DefaultValue = ""
	mask.AddField(f10)

	// name
	f11 := xdommask.NewTextField("name")
	f11.Title = "##name.title##"
	f11.HelpDescription = "##name.help.description##"
	f11.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f11.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f11.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f11.ViewModes = xdommask.DELETE | xdommask.VIEW
	f11.StatusNotNull = "##name.status.notnull##"
	f11.MaxLength = 255
	f11.Size = "400"
	f11.URLVariable = "name"
	f11.DefaultValue = ""
	mask.AddField(f11)

	// description
	f12 := xdommask.NewTextAreaField("description")
	f12.Title = "##description.title##"
	f12.HelpDescription = "##description.help.description##"
	f12.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f12.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f12.ViewModes = xdommask.DELETE | xdommask.VIEW
	f12.MaxLength = 4000
	f12.Width = 400
	f12.Height = 100
	f12.URLVariable = "description"
	f12.DefaultValue = ""
	mask.AddField(f12)

	// Ask what to delete
	f13 := xdommask.NewLOVField("status")
	f13.Title = "##status.title##"
	f13.HelpDescription = "##status.help.description##"
	f13.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f13.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f13.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f13.ViewModes = xdommask.DELETE | xdommask.VIEW
	f13.Options = map[string]string{
		"1": "##status.yes##",
		"2": "##status.no##",
	}
	f13.URLVariable = "status"
	mask.AddField(f13)

	// Ask what to delete
	f21 := xdommask.NewLOVField("prune")
	f21.Title = "##prune.title##"
	f21.HelpDescription = "##prune.help.description##"
	f21.NotNullModes = xdommask.DELETE
	f21.AuthModes = xdommask.DELETE
	f21.HelpModes = xdommask.DELETE
	f21.Options = map[string]string{
		"1": "##prune.yes##",
		"2": "##prune.no##",
	}
	f21.CheckJS = "n = WA.$N('useradmin/profile/editor|" + key + "|profiletoprune'); n.domNodeField.disabled = (value==1); n.checkAll();"
	//	f12.ChangeJS = "function(p) { if (p.id == 'prune') { n = WA.$N('useradmin/profile/editor|" + key + "|grouptoprune'); n.domNodeField.disabled = (p.value==1); n.checkAll(); } }"
	f21.URLVariable = "prune"
	mask.AddField(f21)

	// profile
	f22 := xdommask.NewLOVField("profiletoprune")
	f22.Title = "##profiletoprune.title##"
	f22.HelpDescription = "##profiletoprune.help.description##"
	f22.NotNullModes = xdommask.DELETE
	f22.AuthModes = xdommask.DELETE
	f22.HelpModes = xdommask.DELETE
	f22.Table = user_profile
	f22.Conditions = &xdominion.XConditions{xdominion.NewXCondition("key", "!=", key)}
	f22.Order = &xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	f22.FieldSet = &xdominion.XFieldSet{"key", "name"}
	f22.URLVariable = "profiletoprune"
	mask.AddField(f22)

	// Submit
	f8 := xdommask.NewButtonField("", "submit")
	f8.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE
	f8.TitleInsert = "##formsubmit.insert.title##"
	f8.TitleUpdate = "##formsubmit.update.title##"
	f8.TitleDelete = "##formsubmit.delete.title##"
	mask.AddField(f8)

	// Reset
	f9 := xdommask.NewButtonField("", "reset")
	f9.AuthModes = xdommask.INSERT | xdommask.UPDATE
	f9.TitleInsert = "##formreset.title##"
	f9.TitleUpdate = "##formreset.title##"
	mask.AddField(f9)

	// View
	f91 := xdommask.NewButtonField("", "view")
	f91.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE
	f91.TitleInsert = "##formview.title##"
	f91.TitleUpdate = "##formview.title##"
	f91.TitleDelete = "##formview.title##"
	mask.AddField(f91)

	return nil
}

func predelete(m *xdommask.Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error {

	skey := key.(string)

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)

	prune, _ := rec.GetString("prune")
	if prune == "1" {
		// delete children
		err := userentries.DeleteProfileChildren(ds, skey)
		if err != nil {
			return err
		}
	} else {
		// prune children
		profile, _ := rec.GetString("profiletoprune")
		err := userentries.PruneProfileChildren(ds, skey, profile)
		if err != nil {
			return err
		}
	}
	return nil
}

func Formprofile(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	mask, _ := createMask("formprofile", ctx)
	data, _ := mask.Run(ctx)
	return data
}

func Getaccesses(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	result := []string{}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")

	inaccesses, _ := userentries.GetProfileAccesses(ds, key, 0)
	order := xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	accesses := userentries.GetAccessesList(ds, nil, &order, 0, 0)
	if accesses == nil || len(*accesses) == 0 {
		return "Error"
	}

	// make an array of inprofiles
	ainaccesses := map[string]bool{}
	if inaccesses != nil {
		for _, ia := range *inaccesses {
			ac, _ := ia.GetString("access")
			ainaccesses[ac] = true
		}
	}

	for _, r := range *accesses {
		accesskey, _ := r.GetString("key")
		accessname, _ := r.GetString("name")
		checked := ""
		if ainaccesses[accesskey] {
			checked = " checked=\"checked\""
		}
		result = append(result, "<input type=\"checkbox\""+checked+" onclick=\"useradmin_profile_editor_"+key+"_switchaccess(this, '"+key+"', '"+accesskey+"')\"> <b>"+accesskey+"</b> - "+accessname)
	}
	return strings.Join(result, "<br />")
}

func Getusers(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	result := []string{}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")

	inusers, _ := userentries.GetProfileUsers(ds, key, 0)
	order := xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	users := userentries.GetUsersList(ds, nil, &order, 0, 0)
	if users == nil || len(*users) == 0 {
		return "Error"
	}

	// make an array of inprofiles
	ainusers := map[int]bool{}
	if inusers != nil {
		for _, ip := range *inusers {
			us, _ := ip.GetInt("user")
			ainusers[us] = true
		}
	}

	for _, r := range *users {
		userkey, _ := r.GetInt("key")
		suserkey, _ := r.GetString("key")
		userdata := userentries.GetUserByKey(ds, userkey)
		if userdata == nil {
			return "Error"
		}
		username, _ := userdata.GetString("name")
		checked := ""
		v := ainusers[userkey]
		if v {
			checked = " checked=\"true\""
		}
		result = append(result, "<input type=\"checkbox\""+checked+" onclick=\"useradmin_profile_editor_"+key+"_switchuser(this, '"+key+"', "+suserkey+")\"> <b>"+suserkey+"</b> - "+username)
	}
	return strings.Join(result, "<br />")
}

func Setaccess(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")
	access := ctx.Request.Form.Get("id")
	checked := ctx.Request.Form.Get("checked")

	userentries.SetProfileAccess(ds, key, access, checked == "true")

	return "{\"status\":\"OK\"}"
}

func Setuser(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")
	user := ctx.Request.Form.Get("id")
	iuser, _ := strconv.Atoi(user)
	checked := ctx.Request.Form.Get("checked")

	userentries.SetUserProfile(ds, iuser, key, checked == "true")

	return "{\"status\":\"OK\"}"
}
