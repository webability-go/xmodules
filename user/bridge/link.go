package bridge

import (
	"errors"
	"fmt"
	"plugin"

	"github.com/webability-go/xdominion"
)

var linked bool = false

var GetUsersList func(contextname string) *xdominion.XRecords

func Link(lib *plugin.Plugin) error {
	if linked {
		return nil
	}

	fct, err := lib.Lookup("GetUsersList")
	if err != nil {
		fmt.Println(err)
		return errors.New("ERROR: THE APPLICATION LIBRARY DOES NOT CONTAIN GetUsersList FUNCTION")
	}
	GetUsersList = fct.(func(contextname string) *xdominion.XRecords)

	linked = true
	return nil
}
