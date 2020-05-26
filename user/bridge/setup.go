package bridge

import (
	"errors"
	"fmt"
	"net/http"
	"plugin"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/user/assets"
)

const (
	// DOES NOT MATTER IF THE USER IS OR NOT CONNECTED
	ANY = 1
	// THE USER MUST BE CONNECTED TO USE THE BRIDGE
	USER = 2
)

var ModuleUser *assets.ModuleEntries

// Setup is call ONLY when the app is encapsuled around user directly
func Setup(ctx *serverassets.Context, appname string, connection int) bool {

	// is appname is empty: search for "app" entry in ctx
	if appname == "" {
		appname, _ = ctx.Sysparams.GetString("app")
	}

	// Ask for the plugins we need
	app, ok := ctx.Plugins[appname]
	if !ok {
		// 500 internal error
		http.Error(ctx.Writer, "Library xmodules/app is not available", http.StatusInternalServerError)
		return false
	}

	// Initialize the plugin (just in case)
	err := Link(app)
	if err != nil {
		// 500 internal error
		http.Error(ctx.Writer, "Library xmodules/app could not link with system"+err.Error(), http.StatusInternalServerError)
		return false
	}

	// Verification of session is done during StartContext

	switch connection {
	case ANY:
	case USER:
		sessionid, _ := ctx.Sessionparams.GetString("usersessionid")
		if sessionid == "" {
			http.Error(ctx.Writer, "Error, user not connected", http.StatusUnauthorized)
			return false
		}
	}
	return true
}

var linked = false

func Link(lib *plugin.Plugin) error {
	if linked {
		return nil
	}

	obj, err := lib.Lookup("ModuleUser")
	if err != nil {
		fmt.Println(err)
		fmt.Println(lib)
		return errors.New("ERROR: THE APPLICATION LIBRARY DOES NOT CONTAIN ModuleUser object")
	}
	ok := false
	ModuleUser, ok = obj.(*assets.ModuleEntries)
	if ModuleUser == nil || !ok {
		fmt.Println(lib)
		return errors.New("ERROR: THE APPLICATION LIBRARY DOES NOT CONTAIN ModuleUser valid object")
	}
	linked = true
	return nil
}
