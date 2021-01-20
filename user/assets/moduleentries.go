package assets

import (
	"github.com/webability-go/xdominion"

	//	"github.com/webability-go/xamboo/assets"
	"github.com/webability-go/xmodules/base"
)

type ModuleEntries struct {
	GetAccessesCount func(ds *base.Datasource, cond *xdominion.XConditions) int
	GetAccessesList  func(ds *base.Datasource, cond *xdominion.XConditions, orderby *xdominion.XOrderBy, quantity int, first int) *xdominion.XRecords
}
