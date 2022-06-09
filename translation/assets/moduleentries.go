package assets

import (
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID = "translation"
	VERSION  = "0.0.1"
)

type ModuleEntries struct {
	TranslateOne func(input string) (string, error)
}

func GetEntries(ctx *context.Context) *ModuleEntries {
	me := base.GetEntries(ctx, MODULEID)
	if me == nil {
		return nil
	}
	lme, ok := me.(*ModuleEntries)
	if !ok {
		return nil
	}
	return lme
}
