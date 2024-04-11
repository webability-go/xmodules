package assets

import (
	"log"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID   = "chatgpt"
	VERSION    = "0.0.1"
	DATASOURCE = "chatgptdatasource"

	ACCESSGROUP = "_chatgptadmin"
	ACCESS      = "_chatgptadmin"
)

var Needs = []string{"base"}

type ModuleEntries struct {
	Translate func(ds applications.Datasource, data, lang string) (string, error)
	Ask       func(ds applications.Datasource, data string) (string, error)
}

func GetEntries(logger *log.Logger) *ModuleEntries {
	me := base.GetEntries(logger, MODULEID)
	if me == nil {
		return nil
	}
	lme, ok := me.(*ModuleEntries)
	if !ok {
		return nil
	}
	return lme
}
