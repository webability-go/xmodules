package assets

import (
	"log"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID   = "adminmenu"
	VERSION    = "0.0.1"
	DATASOURCE = "adminmenudatasource"
)

var Needs = []string{"base", "user"}

type ModuleEntries struct {
	AddGroup  func(ds applications.Datasource, key string, name string) error
	GetGroup  func(ds applications.Datasource, key string) (*xdominion.XRecord, error)
	AddOption func(ds applications.Datasource, data *xdominion.XRecord) error
	GetOption func(ds applications.Datasource, key string) (*xdominion.XRecord, error)

	GetMenu func(ds applications.Datasource, group string, father string) (*xdominion.XRecords, error)
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
