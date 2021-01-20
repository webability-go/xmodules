package bridge

import (
	"errors"
	"plugin"

	"github.com/webability-go/xmodules/useradmin/assets"
)

var ModuleUserAdmin *assets.ModuleEntries

var linked = false

func Link(lib *plugin.Plugin) error {
	if linked {
		return nil
	}

	obj, err := lib.Lookup("ModuleUserAdmin")
	if err != nil {
		return errors.New("Error: The application library does not contain the ModuleUserAdmin object")
	}
	ok := false
	ModuleUserAdmin, ok = obj.(*assets.ModuleEntries)
	if ModuleUserAdmin == nil || !ok {
		return errors.New("Error: The application library does not contain a ModuleUserAdmin valid object")
	}
	linked = true
	return nil
}
