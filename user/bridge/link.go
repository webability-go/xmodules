package bridge

import (
	"errors"
	"plugin"

	"github.com/webability-go/xmodules/user/assets"
)

var ModuleUser *assets.ModuleEntries

var linked = false

func Link(lib *plugin.Plugin) error {
	if linked {
		return nil
	}

	obj, err := lib.Lookup("ModuleUser")
	if err != nil {
		return errors.New("Error: The application library does not contain the ModuleUser object")
	}
	ok := false
	ModuleUser, ok = obj.(*assets.ModuleEntries)
	if ModuleUser == nil || !ok {
		return errors.New("Error: The application library does not contain a ModuleUser valid object")
	}
	linked = true
	return nil
}
