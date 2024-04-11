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

	ok := security.Verify(ctx, security.USER, assets.ACCESSLANGUAGE)
	if !ok {
		return ""
	}

	mode := ctx.Request.Form.Get("mode")
	key := ctx.Request.Form.Get("Key")
	if key == "" {
		key = "new"
	}

	params := &xcore.XDataset{
		"FORMSOURCE": createXMLMask("formsource", mode, ctx),
		"KEY":        key,
		"#":          language,
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
	ds := base.TryDatasource(ctx, assets.DATASOURCE)
	if ds == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errordatasource", "build", assets.DATASOURCE))
	}

	kl_traduccionfuente := ds.GetTable("kl_traduccionfuente")
	if kl_traduccionfuente == nil {
		return errors.New(tools.LogMessage(ctx.LoggerError, language, "errortable", "build", "kl_traduccionfuente"))
	}

	//	mode := ctx.Request.Form.Get("mode")
	mask.Table = kl_traduccionfuente
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
	f10 := xdommask.NewIntegerField("clave")
	f10.Title = "##key.title##"
	f10.InRecord = true
	f10.HelpDescription = "##key.help.description##"
	f10.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f10.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f10.ViewModes = xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f10.Auto = true
	f10.AutoMessage = "##key.auto##"
	f10.Size = "400"
	f10.DefaultValue = 0
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

	// configuracion
	f12 := xdommask.NewTextAreaField("config")
	f12.Title = "##config.title##"
	f12.HelpDescription = "##config.help.description##"
	f12.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f12.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f12.ViewModes = xdommask.DELETE | xdommask.VIEW
	f12.MaxLength = 4000
	f12.Width = 400
	f12.Height = 200
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
	f21.CheckJS = "n = WA.$N('translationadmin/source/editor|" + key + "|sourcetoprune'); n.domNodeField.disabled = (value==1); n.checkAll();"
	f21.URLVariable = "prune"
	mask.AddField(f21)

	// source
	f22 := xdommask.NewLOVField("sourcetoprune")
	f22.Title = "##sourcetoprune.title##"
	f22.HelpDescription = "##sourcetoprune.help.description##"
	f22.NotNullModes = xdommask.DELETE
	f22.AuthModes = xdommask.DELETE
	f22.HelpModes = xdommask.DELETE
	f22.Table = kl_traduccionfuente
	f22.Conditions = &xdominion.XConditions{xdominion.NewXCondition("key", "!=", key)}
	f22.Order = &xdominion.XOrder{xdominion.NewXOrderBy("key", xdominion.ASC)}
	f22.FieldSet = &xdominion.XFieldSet{"key", "name"}
	f22.URLVariable = "sourcetoprune"
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

	translationentries := assets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)

	prune, _ := rec.GetString("prune")
	if prune == "1" {
		// delete children
		err := translationentries.DeleteLanguageChildren(ds, skey)
		if err != nil {
			return err
		}
	} else {
		// prune children
		group, _ := rec.GetString("sourcetoprune")
		if len(group) > 1 && group[0] == '_' {
			return errors.New("Error: you cannot paste access children to a system access starting with _")
		}
		err := translationentries.PruneLanguageChildren(ds, skey, group)
		if err != nil {
			return err
		}
	}
	return nil
}

func Formsource(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	mask, _ := createMask("formsource", ctx)
	data, _ := mask.Run(ctx)
	return data
}