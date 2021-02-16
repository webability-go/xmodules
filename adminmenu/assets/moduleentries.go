package assets

import (
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
)

type ModuleEntries struct {
	AddGroup  func(ds applications.Datasource, key string, name string) error
	GetGroup  func(ds applications.Datasource, key string) (*xdominion.XRecord, error)
	AddOption func(ds applications.Datasource, data *xdominion.XRecord) error
	GetOption func(ds applications.Datasource, key string) (*xdominion.XRecord, error)

	GetMenu func(ds applications.Datasource, group string, father string) (*xdominion.XRecords, error)
}
