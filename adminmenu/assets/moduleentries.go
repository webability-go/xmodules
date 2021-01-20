package assets

import (
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/assets"
)

type ModuleEntries struct {
	AddGroup  func(ds assets.Datasource, key string, name string) error
	GetGroup  func(ds assets.Datasource, key string) (*xdominion.XRecord, error)
	AddOption func(ds assets.Datasource, data *xdominion.XRecord) error
	GetOption func(ds assets.Datasource, key string) (*xdominion.XRecord, error)

	GetMenu func(ds assets.Datasource, group string, father string) (*xdominion.XRecords, error)
}
