package bridge

import (
	"errors"
	"plugin"

	"github.com/webability-go/xmodules/base/assets"
)

var ModuleBase *assets.ModuleEntries

var linked = false

func Link(lib *plugin.Plugin) error {
	if linked {
		return nil
	}

	obj, err := lib.Lookup("ModuleBase")
	if err != nil {
		return errors.New("Error: The application library does not contain the ModuleBase object")
	}
	ok := false
	ModuleBase, ok = obj.(*assets.ModuleEntries)
	if ModuleBase == nil || !ok {
		return errors.New("Error: The application library does not contain a ModuleBase valid object")
	}
	linked = true
	return nil
}
