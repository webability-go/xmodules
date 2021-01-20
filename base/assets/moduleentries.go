package assets

import (
	"github.com/webability-go/xamboo/assets"
)

type ModuleEntries struct {
	TryDatasource func(ctx *assets.Context, datasourcename string) assets.Datasource
}
