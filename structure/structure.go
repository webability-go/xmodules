package structure

import (
	"github.com/webability-go/xdominion"

	"xmodules/context"
)

// The structure interface is made to implement a standarized object to use cross modules, graph, memory caches etc.
type Structure interface {

	// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
	ComplementData(sitecontext *context.Context)

	// IsAuthorized returns true if the structure can be used on this site/language/device
	IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool

	// Returns the raw data
	GetData() *xdominion.XRecord
	// Clone the whole structure
	Clone() Structure
}
