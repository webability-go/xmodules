package main

import (
	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/cms/context"

	"xmodules/video/assets"

	"github.com/webability-go/xmodules/user/security"
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

	params := &xcore.XDataset{
		"#": language,
	}

	return template.Execute(params)
}

func Translationdata(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	/*
		key := ctx.Request.Form.Get("Key")
		if key == "" {
			return ""
		}
		ikey, _ := strconv.Atoi(key)
	*/
	result := map[string]interface{}{
		"total":    1,
		"row":      map[int]xdominion.XRecordDef{},
		"subtotal": 0.1,
		"time":     0.1,
		"subtime":  0.1,
	}

	rec := &xdominion.XRecord{
		"key":        1,
		"original":   "abcdef",
		"translated": "ghijkl",
		"commands": "<button class=\"button view\" onclick=\"\">" + language.Get("button.reload") + "</button>" +
			"<button class=\"button view\" onclick=\"\">" + language.Get("button.reload") + "</button>",
	}

	result["row"].(map[int]xdominion.XRecordDef)[0] = rec

	return result
}
