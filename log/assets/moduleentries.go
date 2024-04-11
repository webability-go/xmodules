package assets

import (
	"log"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID    = "log"
	VERSION     = "0.0.1"
	DATASOURCE  = "logdatasource"
	ACCESSGROUP = "_logadmin"
	ACCESS      = "_logadmin"
)

var Needs = []string{"base"}

type ModuleEntries struct {
	GetLogsList  func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	GetLogsCount func(ds applications.Datasource, cond *xdominion.XConditions) int

	AddLog func(ds applications.Datasource, userid int, object, action, keyext string) error
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
