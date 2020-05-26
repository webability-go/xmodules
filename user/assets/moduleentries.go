package assets

import (
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
)

type ModuleEntries struct {
	VerifyUserSession func(ctx *assets.Context, ds *base.Datasource, origin string, device string)
	GetUsersList      func(ds *base.Datasource) *xdominion.XRecords
}
