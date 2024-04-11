package user

import (
	//	"fmt"

	"net/http"
	"regexp"

	//	"github.com/webability-go/xamboo"
	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/tools"
)

// SESSIONS
func VerifyUserSession(ctx *context.Context, ds applications.Datasource, origin string, device string) {

	if !ds.IsModuleAuthorized("user") {
		return
	}

	// security is based on SITE, not DS (not ds.Config)
	config := ctx.Sysparams
	// Any sent session ?
	sessionid := ""
	cookiename, _ := config.GetString("cookiename")
	cookie, _ := ctx.Request.Cookie(cookiename)

	// 1.bis If there is no cookie, there is no session
	if cookie != nil && len(cookie.Value) != 0 {
		sessionid = cookie.Value
	}
	IP := "ip" // ctx.Writer.(*xamboo.CoreWriter).RequestStat.IP

	// verify username, password, OrderSecurity connect/disconnect
	order := ctx.Request.Form.Get("OrderSecurity")

	switch order {
	case "Connect":
		username := ctx.Request.Form.Get("username")
		password := ctx.Request.Form.Get("password")
		// verify against config data
		md5password := tools.GetMD5Hash(password)

		userdata := GetUserByCredentials(ds, username, md5password)
		if userdata != nil {
			// Connect !
			sessionid = CreateSessionUser(ctx, ds, sessionid, IP, origin, device, userdata)
		} else {
			// Disconnect !
			DestroySessionUser(ctx, ds, sessionid)
			return
		}
	case "Disconnect":
		DestroySessionUser(ctx, ds, sessionid)
		return
	}

	if sessionid == "" { // there is no session
		return
	}
	sessiondata := GetSession(ds, sessionid)
	if sessiondata == nil {
		return
	}

	checkIP, _ := config.GetBool("checkip")
	sessionip, _ := sessiondata.GetString("ip")
	if checkIP && IP != sessionip {
		DestroySessionUser(ctx, ds, sessionid)
		return
	}

	// set user data, update session

	// link session with ctx
	//	ctx.Sessionparams.Set("sessionid", sessionid)
	//	ctx.Sessionparams.Set("clientid", clientid)

	userkey, _ := sessiondata.GetInt("user")
	userdata := GetUserByKey(ds, userkey)
	username, _ := userdata.GetString("name")

	ctx.Sessionparams.Set("usersessionid", sessionid)
	ctx.Sessionparams.Set("userkey", userkey)
	ctx.Sessionparams.Set("username", username)
	ctx.Sessionparams.Set("usersession", sessiondata)
	ctx.Sessionparams.Set("userdata", userdata)
}

func CreateSessionUser(ctx *context.Context, ds applications.Datasource, sessionid string, IP string, origin string, device string, user *StructureUser) string {

	config := ctx.Sysparams

	match, _ := regexp.MatchString("[a-zA-Z0-9]{24}", sessionid)
	if !match {
		sessionid = tools.UUID()
	}

	userkey, _ := user.Data.GetInt("key")
	sessionid = CreateSession(ds, userkey, sessionid, IP, origin, device)
	if sessionid == "" {
		return ""
	}

	cookiedomain, _ := config.GetString("cookiedomain")
	cookiename, _ := config.GetString("cookiename")
	http.SetCookie(ctx.Writer, &http.Cookie{Domain: cookiedomain, Name: cookiename, Value: sessionid, Path: "/"})
	return sessionid
}

func DestroySessionUser(ctx *context.Context, ds applications.Datasource, sessionid string) {

	config := ctx.Sysparams
	cookiedomain, _ := config.GetString("cookiedomain")
	cookiename, _ := config.GetString("cookiename")

	http.SetCookie(ctx.Writer, &http.Cookie{Domain: cookiedomain, Name: cookiename, Value: "", Path: "/", MaxAge: -1})

	// destroys the session in DB
	CloseSession(ds, sessionid)
}

// SECURITY Access
func HasAccess(ds applications.Datasource, userid int, args ...interface{}) bool {

	// 1. check direct acceses
	// 2. or check if in profile
	userdata := GetUserByKey(ds, userid)
	if userdata == nil {
		return false
	}

	// superuser is always all access
	super, _ := userdata.GetString("status")
	if super == "S" {
		return true
	}
	if super == "X" { // desactivated user, no access
		return false
	}

	access := ""
	extendedaccess := ""
	ok := false
	// any args ?
	if len(args) > 0 {
		access, ok = args[0].(string)
		if !ok {
			//			http.Error(ctx.Writer, "Error on user access verification", http.StatusUnauthorized)
			return false
		}
	}
	if len(args) > 1 {
		extendedaccess = args[1].(string)
		if !ok {
			//			http.Error(ctx.Writer, "Error on user extended access verification", http.StatusUnauthorized)
			return false
		}
	}

	if extendedaccess == "" {
		// simple access:
		// 1. check into user_useraccess
		acc, err := GetUserAccess(ds, userid, access)
		if err != nil {
			ds.Log("xmodules::user::HasAccess:1:" + err.Error())
			return false
		}
		if acc != nil { // if any record, direct right on access
			denied, _ := acc.GetInt("denied")
			return denied == 0 // if denied == 1: do not acces, else, access
		}
		// If we are here, the access is herency of profile
		hasaccess, err := HasProfilesAccessUser(ds, userid, access)
		if err != nil {
			ds.Log("xmodules::user::HasAccess:2:" + err.Error())
			return false
		}
		return hasaccess
	} else {

	}
	/*
		// is access into user accesses ?
		acc := GetUserAccess(ds, userid, access)

		// is access in profile ?
		pr := GetUserProfiles(ds, userid)
		// loop sobre progiles
		accpr := GetProfileAccess(ds, userid, access)
	*/

	return false
}

/*
public function hasAccess($claveusuario, $derecho)
{
	if (!$claveusuario)
		return false;
	$datausuario = $this->getUsuarioDataByKey($claveusuario);
	if (!$datausuario)
		return false;
	if ($datausuario->estatus == "S") // super usuario
		return true;

	// 1. check access or extended
	$data = $this->kl_adminderechousuario->doSelectCondition(array(new \dominion\DB_Condition('usuario', '=', $claveusuario), new \dominion\DB_Condition('derecho', '=', $derecho, 'and')));

	if ($data && $data[0]->estatus == 1) // Authorized
		return true;           // if denied for sure we return false
	if ($data && $data[0]->estatus == 2) // Denied
		return false;

	$profiles = $this->kl_adminperfilusuario->doSelectCondition(array(new \dominion\DB_Condition('usuario', '=', $claveusuario)));
	if ($profiles)
	{
		foreach($profiles as $pr)
		{
			$prof = $this->kl_adminperfilderecho->doSelectCondition(array(new \dominion\DB_Condition('perfil', '=', $pr->perfil), new \dominion\DB_Condition('derecho', '=', $derecho, 'and')));
			if ($prof)
				return true;    // found: authorized
		}
	}
	return false;
}
*/
