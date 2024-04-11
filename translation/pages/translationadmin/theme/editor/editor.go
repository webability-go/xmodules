package main

import (
	"encoding/xml"
	"errors"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xdommask"

	//	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"xmodules/translation/assets"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user/security"
)

var language *xcore.XLanguage

func Run(ctx *context.Context, template *xcore.XTemplate, xlanguage *xcore.XLanguage, e interface{}) interface{} {

	if language == nil {
		language = xlanguage
	}

	ok := security.Verify(ctx, security.USER, assets.ACCESSTHEME)
	if !ok {
		return ""
	}

	mode := ctx.Request.Form.Get("mode")
	key := ctx.Request.Form.Get("Key")
	if key == "" {
		key = "new"
	}

	params := &xcore.XDataset{
		"FORMTHEME": createXMLMask("formtheme", mode, ctx),
		"KEY":       key,
		"#":         language,
	}

	return template.Execute(params)
}

func createMask(id string, ctx *context.Context) (*xdommask.Mask, error) {

	hooks := xdommask.MaskHooks{
		Build:     build,
		PreInsert: preinsert,
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

	kl_traducciontema := ds.GetTable("kl_traducciontema")
	if kl_traducciontema == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errortable", "build", "kl_traducciontema"))
	}

	//	mode := ctx.Request.Form.Get("mode")
	mask.Table = kl_traducciontema
	key := ctx.Request.Form.Get(mask.VarKey)
	if key != "" {
		mask.Key = key
	}

	mask.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	mask.KeyField = "clave"

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
	f10 := xdommask.NewTextField("clave")
	f10.Title = "##key.title##"
	f10.InRecord = true
	f10.HelpDescription = "##key.help.description##"
	f10.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f10.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f10.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f10.ViewModes = xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f10.StatusNotNull = "##key.status.notnull##"
	f10.StatusBadFormat = "##key.status.badformat##"
	f10.MaxLength = 5
	f10.Size = "400"
	f10.Format = "^[a-z][a-z0-9-_]{1,5}$"
	f10.FormatJS = "^[a-z][a-z0-9-_]{1,5}$"
	f10.DefaultValue = ""
	mask.AddField(f10)

	// name
	f11 := xdommask.NewTextField("nombre")
	f11.Title = "##name.title##"
	f11.HelpDescription = "##name.help.description##"
	f11.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f11.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f11.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f11.ViewModes = xdommask.DELETE | xdommask.VIEW
	f11.StatusNotNull = "##name.status.notnull##"
	f11.MaxLength = 200
	f11.Size = "400"
	mask.AddField(f11)

	// translation status
	f111 := xdommask.NewLOVField("estatus")
	f111.Title = "##status.title##"
	f111.InRecord = true
	f111.HelpDescription = "##status.help.description##"
	f111.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f111.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f111.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f111.ViewModes = xdommask.DELETE | xdommask.VIEW
	f111.StatusNotNull = "##status.status.notnull##"
	f111.Options = map[string]string{
		"1": "1 / ##status.active##",
		"2": "2 / ##status.down##",
	}
	f111.DefaultValue = "1"
	mask.AddField(f111)

	// description
	f12 := xdommask.NewTextAreaField("descripcion")
	f12.Title = "##description.title##"
	f12.HelpDescription = "##description.help.description##"
	f12.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f12.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f12.ViewModes = xdommask.DELETE | xdommask.VIEW
	f12.MaxLength = 4000
	f12.Width = 400
	f12.Height = 100
	mask.AddField(f12)

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
	f21.CheckJS = "n = WA.$N('translationadmin/theme/editor|" + key + "|themetoprune'); n.domNodeField.disabled = (value==1); n.checkAll();"
	f21.URLVariable = "prune"
	mask.AddField(f21)

	// theme
	f22 := xdommask.NewLOVField("themetoprune")
	f22.Title = "##themetoprune.title##"
	f22.HelpDescription = "##themetoprune.help.description##"
	f22.NotNullModes = xdommask.DELETE
	f22.AuthModes = xdommask.DELETE
	f22.HelpModes = xdommask.DELETE
	f22.Table = kl_traducciontema
	f22.Conditions = &xdominion.XConditions{xdominion.NewXCondition("clave", "!=", key), xdominion.NewXCondition("clave", "not like", "\\_%", "and")}
	f22.Order = &xdominion.XOrder{xdominion.NewXOrderBy("clave", xdominion.ASC)}
	f22.FieldSet = &xdominion.XFieldSet{"clave", "nombre"}
	f22.URLVariable = "themetoprune"
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

	// assign ID as max + 1
	max, err := m.Table.Max("orden")
	if err == nil {
		rec.Set("orden", max.(int64)+1)
	}
	key, _ := rec.GetString("clave")
	rec.Set("label", key)

	return nil
}

func predelete(m *xdommask.Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error {
	skey := key.(string)

	translationentries := assets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)

	prune, _ := rec.GetString("prune")
	if prune == "1" {
		// delete children
		err := translationentries.DeleteThemeChildren(ds, skey)
		if err != nil {
			return err
		}
	} else {
		// prune children
		group, _ := rec.GetString("themetoprune")
		if len(group) > 1 && group[0] == '_' {
			return errors.New("Error: you cannot paste access children to a system access starting with _")
		}
		err := translationentries.PruneThemeChildren(ds, skey, group)
		if err != nil {
			return err
		}
	}
	return nil
}

func Formtheme(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	mask, _ := createMask("formtheme", ctx)
	data, _ := mask.Run(ctx)
	return data
}
