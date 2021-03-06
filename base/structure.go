package base

import (
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
)

// Structure interface is made to implement a standarized object to use cross modules, graph, memory caches etc.
type Structure interface {

	// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
	ComplementData(ds applications.Datasource)

	// IsAuthorized returns true if the structure can be used on this site/language/device
	IsAuthorized(ds applications.Datasource, site string, language string, device string) bool

	// GetData Returns the raw data
	GetData() *xdominion.XRecord

	// Clone will clone the whole structure
	Clone() Structure
}
