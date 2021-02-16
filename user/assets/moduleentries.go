package assets

import (
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
)

type ModuleEntries struct {
	// Accesses
	GetAccessesCount func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetAccessesList  func(ds applications.Datasource, cond *xdominion.XConditions, orderby *xdominion.XOrderBy, quantity int, first int) *xdominion.XRecords

	// User Params
	SetUserParam func(ds applications.Datasource, user int, param string, value interface{})
	AddUserParam func(ds applications.Datasource, user int, param string, value interface{})
	GetUserParam func(ds applications.Datasource, user int, param string) string
	DelUserParam func(ds applications.Datasource, user int, param string)

	HasAccess func(ds applications.Datasource, user int, access string, extra string) bool
}
