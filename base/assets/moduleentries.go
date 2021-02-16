package assets

import (
	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"
)

type ModuleEntries struct {
	TryDatasource func(ctx *context.Context, datasourcename string) applications.Datasource
}
