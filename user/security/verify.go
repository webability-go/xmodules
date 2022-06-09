package security

import (
	"net/http"

	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/user/assets"
)

const (
	// DOES NOT MATTER IF THE USER IS OR NOT CONNECTED
	ANY = 1
	// THE USER MUST BE CONNECTED TO USE THE LIBARY
	USER = 2
	// No access needed
	FREEACCESS = ""
)

// Verify connected user, verify security accesses
func Verify(ctx *context.Context, connection int, args ...interface{}) bool {

	// is appname is empty: search for "app" entry in ctx
	appname, _ := ctx.Sysparams.GetString("adminapp")
	if appname == "" {
		http.Error(ctx.Writer, "Admin Library name is not available in config file (parameter adminapp missing)", http.StatusInternalServerError)
		return false
	}

	// Verification of session is done during StartContext of application, and set the Sessionparam
	switch connection {
	case ANY:
		return true
	case USER:
		sessionid, _ := ctx.Sessionparams.GetString("usersessionid")
		if sessionid == "" {
			http.Error(ctx.Writer, "Error, user not connected", http.StatusForbidden)
			return false
		}
	}

	access := ""
	ok := false
	// any args ?
	if len(args) > 0 {
		access, ok = args[0].(string)
		if !ok {
			http.Error(ctx.Writer, "Error on user access verification", http.StatusInternalServerError)
			return false
		}
	}

	// security access
	if access != "" {
		// Check security
		acc := HasAccess(ctx, args...)
		if !acc {
			http.Error(ctx.Writer, "Error, user not authorized", http.StatusForbidden)
			return false
		}
	}

	return true
}

// SECURITY Access
func HasAccess(ctx *context.Context, args ...interface{}) bool {

	userentries := assets.GetEntries(ctx)
	if userentries == nil {
		http.Error(ctx.Writer, "Error, no security module", http.StatusInternalServerError)
		return false
	}

	ds := base.TryDatasource(ctx, assets.DATASOURCE)
	if ds == nil {
		http.Error(ctx.Writer, "Error, no security datasource", http.StatusInternalServerError)
		return false
	}
	userkey, _ := ctx.Sessionparams.GetInt("userkey")

	return userentries.HasAccess(ds, userkey, args...)
}
