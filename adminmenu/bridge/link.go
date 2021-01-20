package bridge

import (
	"errors"
	"plugin"

	"github.com/webability-go/xmodules/adminmenu/assets"
)

var ModuleAdminMenu *assets.ModuleEntries

var linked = false

func Link(lib *plugin.Plugin) error {
	if linked {
		return nil
	}

	obj, err := lib.Lookup("ModuleAdminMenu")
	if err != nil {
		return errors.New("Error: The application library does not contain the ModuleUserAdmin object")
	}
	ok := false
	ModuleAdminMenu, ok = obj.(*assets.ModuleEntries)
	if ModuleAdminMenu == nil || !ok {
		return errors.New("Error: The application library does not contain a ModuleAdminMenu valid object")
	}
	linked = true
	return nil
}
