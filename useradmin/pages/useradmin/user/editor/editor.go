package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
		"FORMUSER": createXMLMask("formuser", mode, ctx),
		"KEY":      key,
		"#":        language,
	}

	return template.Execute(params)
}

func createMask(id string, ctx *context.Context) (*xdommask.Mask, error) {

	hooks := xdommask.MaskHooks{
		Build:     build,
		PreInsert: preinsert,
		PreUpdate: preupdate,
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
	ds := base.TryDatasource(ctx, assets.DATASOURCE)
	if ds == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errordatasource", "build", assets.DATASOURCE))
	}

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errortable", "build", "user_user"))
	}

	mask.Table = user_user
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
	f10 := xdommask.NewIntegerField("key")
	f10.Title = "##key.title##"
	f10.HelpDescription = "##key.help.description##"
	f10.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f10.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f10.ViewModes = xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f10.Auto = true
	f10.AutoMessage = "##key.auto##"
	f10.Size = "400"
	f10.URLVariable = "key"
	f10.DefaultValue = 0
	mask.AddField(f10)

	// user status
	f11 := xdommask.NewLOVField("status")
	f11.Title = "##status.title##"
	f11.HelpDescription = "##status.help.description##"
	f11.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f11.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f11.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f11.ViewModes = xdommask.DELETE | xdommask.VIEW
	f11.StatusNotNull = "##status.status.notnull##"
	f11.Options = map[string]string{
		"A": "A / ##status.active##",
		"S": "S / ##status.superuser##",
		"X": "X / ##status.down##",
	}
	f11.URLVariable = "status"
	f11.DefaultValue = "A"
	mask.AddField(f11)

	// name
	f12 := xdommask.NewTextField("name")
	f12.Title = "##name.title##"
	f12.HelpDescription = "##name.help.description##"
	f12.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f12.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f12.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f12.ViewModes = xdommask.DELETE | xdommask.VIEW
	f12.StatusNotNull = "##name.status.notnull##"
	f12.MaxLength = 200
	f12.Size = "400"
	f12.URLVariable = "name"
	f12.DefaultValue = ""
	mask.AddField(f12)

	// username
	f13 := xdommask.NewTextField("un")
	f13.Title = "##un.title##"
	f13.HelpDescription = "##un.help.description##"
	f13.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f13.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f13.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f13.ViewModes = xdommask.DELETE | xdommask.VIEW
	f13.StatusNotNull = "##un.status.notnull##"
	f13.MaxLength = 200
	f13.Size = "400"
	f13.URLVariable = "un"
	f13.DefaultValue = ""
	mask.AddField(f13)

	// password
	f14 := xdommask.NewMaskedField("pw")
	f14.Title = "##pw.title##"
	f14.HelpDescription = "##pw.help.description##"
	//	f14.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f14.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f14.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f14.ViewModes = xdommask.DELETE | xdommask.VIEW
	//	f14.StatusNotNull = "##pw.status.notnull##"
	f14.MD5Encrypted = true
	f14.MaxLength = 200
	f14.Size = "400"
	f14.URLVariable = "pw"
	f14.DefaultValue = ""
	mask.AddField(f14)

	// mail
	f15 := xdommask.NewTextField("mail")
	f15.Title = "##mail.title##"
	f15.HelpDescription = "##mail.help.description##"
	f15.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f15.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f15.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f15.ViewModes = xdommask.DELETE | xdommask.VIEW
	f15.StatusNotNull = "##mail.status.notnull##"
	f15.MaxLength = 255
	f15.Size = "400"
	f15.URLVariable = "mail"
	f15.DefaultValue = ""
	mask.AddField(f15)

	// gender
	f16 := xdommask.NewLOVField("sex")
	f16.Title = "##sex.title##"
	f16.HelpDescription = "##sex.help.description##"
	f16.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f16.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f16.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f16.ViewModes = xdommask.DELETE | xdommask.VIEW
	f16.StatusNotNull = "##sex.status.notnull##"
	f16.Options = map[string]string{
		"M": "##sex.masc##",
		"F": "##sex.fem##",
	}
	f16.URLVariable = "sex"
	f16.DefaultValue = "M"
	mask.AddField(f16)

	// fields
	f17 := xdommask.NewTextAreaField("fields")
	f17.Title = "##fields.title##"
	f17.HelpDescription = "##fields.help.description##"
	f17.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f17.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f17.ViewModes = xdommask.DELETE | xdommask.VIEW
	f17.MaxLength = 4000
	f17.Width = 400
	f17.Height = 150
	f17.URLVariable = "fields"
	f17.DefaultValue = ""
	mask.AddField(f17)

	// info
	f18 := xdommask.NewTextAreaField("info")
	f18.Title = "##info.title##"
	f18.HelpDescription = "##info.help.description##"
	f18.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f18.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f18.ViewModes = xdommask.DELETE | xdommask.VIEW
	f18.MaxLength = 4000
	f18.Width = 400
	f18.Height = 150
	f18.URLVariable = "info"
	f18.DefaultValue = ""
	mask.AddField(f18)

	// father, responsible
	f19 := xdommask.NewLOVField("father")
	f19.Title = "##father.title##"
	f19.HelpDescription = "##father.help.description##"
	f19.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f19.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f19.ViewModes = xdommask.DELETE | xdommask.VIEW
	f19.Table = user_user
	if key != "" {
		f19.Conditions = &xdominion.XConditions{xdominion.NewXCondition("key", "!=", key)}
	}
	f19.Order = &xdominion.XOrder{xdominion.NewXOrderBy("name", xdominion.ASC)}
	f19.FieldSet = &xdominion.XFieldSet{"key", "name"}
	f19.NullOnEmpty = true
	f19.URLVariable = "father"
	mask.AddField(f19)

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

func preinsert(m *xdommask.Mask, ctx *context.Context, rec *xdominion.XRecord) error {
	rec.Set("creationdate", time.Now())
	rec.Set("lastmodif", time.Now())
	return nil
}

func preupdate(m *xdommask.Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) error {
	newrec.Set("lastmodif", time.Now())
	return nil
}

func predelete(m *xdommask.Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error {
	ikey := key.(int)

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)

	// delete children
	err := userentries.DeleteUserChildren(ds, ikey)
	return err
}

func Formuser(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	mask, _ := createMask("formuser", ctx)
	data, _ := mask.Run(ctx)
	return data
}

func Getprofiles(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	result := []string{}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")
	ikey, _ := strconv.Atoi(key)

	inprofiles, err := userentries.GetUserProfiles(ds, ikey, 10)
	fmt.Println("IN PROFILES:", key, ikey, inprofiles, err)
	order := xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	profiles := userentries.GetProfilesList(ds, nil, &order, 0, 0)
	if profiles == nil || len(*profiles) == 0 {
		return "Error"
	}

	// make an array of inprofiles
	ainprofiles := map[string]bool{}
	if inprofiles != nil {
		for _, ip := range *inprofiles {
			pr, _ := ip.GetString("profile")
			ainprofiles[pr] = true
		}
	}
	fmt.Println("A IN PROFILES:", ainprofiles)
	for _, r := range *profiles {
		profilekey, _ := r.GetString("key")
		profilename, _ := r.GetString("name")
		checked := ""
		if ainprofiles[profilekey] {
			checked = " checked=\"checked\""
		}
		result = append(result, "<input type=\"checkbox\""+checked+" onclick=\"useradmin_user_editor_"+key+"_switchprofile(this, '"+key+"', '"+profilekey+"')\"> <b>"+profilekey+"</b> - "+profilename)
	}
	return strings.Join(result, "<br />")
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
	ikey, _ := strconv.Atoi(key)

	inusers, _ := userentries.GetUserAccesses(ds, ikey, 0)
	order := xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	accesses := userentries.GetAccessesList(ds, nil, &order, 0, 0)
	if accesses == nil || len(*accesses) == 0 {
		return "Error"
	}

	// make an array of inprofiles
	ainusers := map[string]int{}
	if inusers != nil {
		for _, ip := range *inusers {
			ac, _ := ip.GetString("access")
			denied, _ := ip.GetInt("denied")
			ainusers[ac] = denied
		}
	}

	bydefault := language.Get("accesses.hierarchy")
	authorized := language.Get("accesses.yes")
	denied := language.Get("accesses.no")
	for _, r := range *accesses {
		accesskey, _ := r.GetString("key")
		accessdata := userentries.GetAccessByKey(ds, accesskey)
		if accessdata == nil {
			return "Error"
		}
		accessname, _ := accessdata.GetString("name")
		checked1 := ""
		checked2 := ""
		checked3 := ""
		v, ok := ainusers[accesskey]
		if !ok {
			checked1 = " checked=\"true\""
		} else if v == 0 {
			checked2 = " checked=\"true\""
		} else { // v == 1
			checked3 = " checked=\"true\""
		}
		result = append(result, "<input type=\"radio\""+checked1+" name=\""+accesskey+"\" onclick=\"useradmin_user_editor_"+key+"_switchaccess(this, '"+key+"', '"+accesskey+"', -1)\" value=\"-1\"> "+bydefault+
			"<input type=\"radio\""+checked2+" name=\""+accesskey+"\" onclick=\"useradmin_user_editor_"+key+"_switchaccess(this, '"+key+"', '"+accesskey+"', 0)\" value=\"0\"> "+authorized+
			"<input type=\"radio\""+checked3+" name=\""+accesskey+"\" onclick=\"useradmin_user_editor_"+key+"_switchaccess(this, '"+key+"', '"+accesskey+"', 1)\" value=\"1\"> "+denied+" | "+
			" <b>"+accesskey+"</b> - "+accessname)
	}
	return strings.Join(result, "<br />")
}

func Setprofile(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")
	ikey, _ := strconv.Atoi(key)
	profile := ctx.Request.Form.Get("id")
	checked := ctx.Request.Form.Get("checked")

	userentries.SetUserProfile(ds, ikey, profile, checked == "true")

	return "{\"status\":\"OK\"}"
}

func Setaccess(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	user := ctx.Request.Form.Get("key")
	iuser, _ := strconv.Atoi(user)
	access := ctx.Request.Form.Get("id")
	val := ctx.Request.Form.Get("val")
	ival, _ := strconv.Atoi(val)

	userentries.SetUserAccess(ds, iuser, access, ival)

	return "{\"status\":\"OK\"}"
}
