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
		"FORMACCESS": createXMLMask("formaccess", mode, ctx),
		"KEY":        key,
		"#":          language,
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
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	if ds == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errordatasource", "build", userassets.DATASOURCE))
	}

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errortable", "build", "user_access"))
	}

	user_accessgroup := ds.GetTable("user_accessgroup")
	if user_accessgroup == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errortable", "build", "user_accessgroup"))
	}

	mode := ctx.Request.Form.Get("mode")
	mask.Table = user_access
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
	f10.StatusBadFormat = "##key.status.badformat##"
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

	// group
	f13 := xdommask.NewLOVField("group")
	f13.Title = "##group.title##"
	f13.HelpDescription = "##group.help.description##"
	f13.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f13.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f13.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f13.ViewModes = xdommask.DELETE | xdommask.VIEW
	f13.Table = user_accessgroup
	if mode != "4" {
		f13.Conditions = &xdominion.XConditions{xdominion.NewXCondition("key", "not like", "\\_%")}
	}
	f13.Order = &xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	f13.FieldSet = &xdominion.XFieldSet{"key", "name"}
	f13.URLVariable = "group"
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
	f21.CheckJS = "n = WA.$N('useradmin/access/editor|" + key + "|accesstoprune'); n.domNodeField.disabled = (value==1); n.checkAll();"
	//	f12.ChangeJS = "function(p) { if (p.id == 'prune') { n = WA.$N('useradmin/accessgroup/editor|" + key + "|grouptoprune'); n.domNodeField.disabled = (p.value==1); n.checkAll(); } }"
	f21.URLVariable = "prune"
	mask.AddField(f21)

	// group
	f22 := xdommask.NewLOVField("accesstoprune")
	f22.Title = "##accesstoprune.title##"
	f22.HelpDescription = "##accesstoprune.help.description##"
	f22.NotNullModes = xdommask.DELETE
	f22.AuthModes = xdommask.DELETE
	f22.HelpModes = xdommask.DELETE
	f22.Table = user_access
	f22.Conditions = &xdominion.XConditions{xdominion.NewXCondition("key", "!=", key), xdominion.NewXCondition("key", "not like", "\\_%", "and")}
	f22.Order = &xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	f22.FieldSet = &xdominion.XFieldSet{"key", "name"}
	f22.URLVariable = "accesstoprune"
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

func preinsert(m *xdommask.Mask, ctx *context.Context, rec *xdominion.XRecord) error {
	key, _ := rec.GetString("key")
	if len(key) > 1 && key[0] == '_' {
		return errors.New("Error: you cannot insert a key starting with _")
	}
	return nil
}

func preupdate(m *xdommask.Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) error {
	skey := key.(string)
	if len(skey) > 1 && skey[0] == '_' {
		return errors.New("Error: you cannot update a key starting with _")
	}
	newkey, _ := newrec.GetString("key")
	if len(newkey) > 1 && newkey[0] == '_' {
		return errors.New("Error: you cannot update a key starting with _")
	}
	return nil
}

func predelete(m *xdommask.Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error {
	skey := key.(string)
	if len(skey) > 1 && skey[0] == '_' {
		return errors.New("Error: you cannot delete a key starting with _")
	}

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)

	prune, _ := rec.GetString("prune")
	if prune == "1" {
		// delete children
		err := userentries.DeleteAccessChildren(ds, skey)
		if err != nil {
			return err
		}
	} else {
		// prune children
		group, _ := rec.GetString("accesstoprune")
		if len(group) > 1 && group[0] == '_' {
			return errors.New("Error: you cannot paste access children to a system access starting with _")
		}
		err := userentries.PruneAccessChildren(ds, skey, group)
		if err != nil {
			return err
		}
	}
	return nil
}

func Formaccess(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	mask, _ := createMask("formaccess", ctx)
	data, _ := mask.Run(ctx)
	return data
}

func Getprofiles(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	result := []string{}

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")

	inprofiles, _ := userentries.GetAccessProfiles(ds, key, 0)
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

	for _, r := range *profiles {
		profilekey, _ := r.GetString("key")
		profilename, _ := r.GetString("name")
		checked := ""
		if ainprofiles[profilekey] {
			checked = " checked=\"checked\""
		}
		result = append(result, "<input type=\"checkbox\""+checked+" onclick=\"useradmin_access_editor_"+key+"_switchprofile(this, '"+key+"', '"+profilekey+"')\"> <b>"+profilekey+"</b> - "+profilename)
	}
	return strings.Join(result, "<br />")
}

func Getusers(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	result := []string{}

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")

	inusers, _ := userentries.GetAccessUsers(ds, key, 0)
	order := xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	users := userentries.GetUsersList(ds, nil, &order, 0, 0)
	if users == nil || len(*users) == 0 {
		return "Error"
	}

	// make an array of inprofiles
	ainusers := map[int]int{}
	if inusers != nil {
		for _, ip := range *inusers {
			us, _ := ip.GetInt("user")
			denied, _ := ip.GetInt("denied")
			ainusers[us] = denied
		}
	}

	bydefault := language.Get("users.hierarchy")
	authorized := language.Get("users.yes")
	denied := language.Get("users.no")
	for _, r := range *users {
		userkey, _ := r.GetInt("key")
		suserkey, _ := r.GetString("key")
		userdata := userentries.GetUserByKey(ds, userkey)
		if userdata == nil {
			return "Error"
		}
		username, _ := userdata.GetString("name")
		checked1 := ""
		checked2 := ""
		checked3 := ""
		v, ok := ainusers[userkey]
		if !ok {
			checked1 = " checked=\"true\""
		} else if v == 0 {
			checked2 = " checked=\"true\""
		} else { // v == 1
			checked3 = " checked=\"true\""
		}
		result = append(result, "<input type=\"radio\""+checked1+" name=\""+suserkey+"\" onclick=\"useradmin_access_editor_"+key+"_switchuser(this, '"+key+"', '"+suserkey+"', -1)\" value=\"-1\"> "+bydefault+
			"<input type=\"radio\""+checked2+" name=\""+suserkey+"\" onclick=\"useradmin_access_editor_"+key+"_switchuser(this, '"+key+"', '"+suserkey+"', 0)\" value=\"0\"> "+authorized+
			"<input type=\"radio\""+checked3+" name=\""+suserkey+"\" onclick=\"useradmin_access_editor_"+key+"_switchuser(this, '"+key+"', '"+suserkey+"', 1)\" value=\"1\"> "+denied+" | "+
			" <b>"+suserkey+"</b> - "+username)
	}
	return strings.Join(result, "<br />")
}

func Setprofile(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")
	profile := ctx.Request.Form.Get("id")
	checked := ctx.Request.Form.Get("checked")

	userentries.SetProfileAccess(ds, profile, key, checked == "true")

	return "{\"status\":\"OK\"}"
}

func Setuser(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	userentries := userassets.GetEntries(ctx)
	ds := base.TryDatasource(ctx, userassets.DATASOURCE)
	key := ctx.Request.Form.Get("key")
	user := ctx.Request.Form.Get("id")
	iuser, _ := strconv.Atoi(user)
	val := ctx.Request.Form.Get("val")
	ival, _ := strconv.Atoi(val)

	userentries.SetUserAccess(ds, iuser, key, ival)

	return "{\"status\":\"OK\"}"
}
