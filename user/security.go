package user

import (
	"net/http"
	"regexp"

	"github.com/webability-go/xamboo"
	"github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
)

func VerifyUserSession(ctx *assets.Context, ds *base.Datasource, origin string, device string) {

	if !ds.IsModuleAuthorized("user") {
		return
	}

	config := ds.Config
	// Any sent session ?
	sessionid := ""
	cookiename, _ := config.GetString("cookiename")
	cookie, _ := ctx.Request.Cookie(cookiename)
	// 1.bis If there is no cookie, there is no session
	if cookie != nil && len(cookie.Value) != 0 {
		sessionid = cookie.Value
	}
	IP := ctx.Writer.(*xamboo.CoreWriter).RequestStat.IP

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
	userdata := GetUser(ds, userkey)
	username, _ := userdata.Data.GetString("name")

	ctx.Sessionparams.Set("usersessionid", sessionid)
	ctx.Sessionparams.Set("userkey", userkey)
	ctx.Sessionparams.Set("username", username)
	ctx.Sessionparams.Set("usersession", sessiondata)
	ctx.Sessionparams.Set("userdata", userdata.Data)
}

func CreateSessionUser(ctx *assets.Context, ds *base.Datasource, sessionid string, IP string, origin string, device string, user *StructureUser) string {

	config := ds.Config

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

func DestroySessionUser(ctx *assets.Context, ds *base.Datasource, sessionid string) {

	config := ds.Config
	cookiedomain, _ := config.GetString("cookiedomain")
	cookiename, _ := config.GetString("cookiename")

	http.SetCookie(ctx.Writer, &http.Cookie{Domain: cookiedomain, Name: cookiename, Value: "", Path: "/", MaxAge: -1})

	// destroys the session in DB
	CloseSession(ds, sessionid)
}

/* Verify cookie session against DB
func LinkSessionUser(ctx *assets.Context) {

	sitecontextname, _ := ctx.Sysparams.GetString("sitecontext")
	sitecontext := base.Sites.GetContext(sitecontextname)
	if sitecontext == nil {
		return
	}
	// kiwi-gr, kiwi-im, crafto-gr, crafto-im are authorized normally. control, kiwi7, crafto7, central, cdn, are not
	if !context.IsModuleAuthorized(sitecontext, "client") {
		return
	}

	var clientsession *xdominion.XRecord

	sessionid := ctx.Request.Form.Get("token")
	if sessionid != "" {
		clientsession = client.GetSession(sessionid)
	}

	config := base.Sites.GetContext("kiwi-gr").Config

	if clientsession == nil {
		// 1. check the cookie session
		cookiename, _ := config.GetString("cookiename")
		cookie, _ := ctx.Request.Cookie(cookiename)
		// 1.bis If there is no cookie, there is no session
		if cookie == nil || len(cookie.Value) == 0 {
			return
		}
		sessionid = cookie.Value
		clientsession = client.GetSession(sessionid)
	}

	if clientsession == nil {
		// Si no hay sesion repertoriada, destruye todo
		DestroySessionClient(ctx, sessionid)
		return
	}

	checkIP, _ := ctx.Sysparams.GetBool("checkip")
	IP, _ := clientsession.GetString("ip")
	IPClient := ctx.Writer.(*server.CoreWriter).RequestStat.IP
	if checkIP && IP != IPClient {
		DestroySessionClient(ctx, sessionid)
		return
	}

	// Actualiza el fin de la sesion con el tiempo actual
	client.SetSessionTime(sessionid)

	clave, _ := clientsession.GetInt("chef")
	if clave != 0 {
		client.SetClientTime(clave)
		clientdata, _ := graph.GetClient(sitecontext, clave)
		if clientdata != nil {
			ctx.Sessionparams.Set("sessionid", sessionid)
			ctx.Sessionparams.Set("clientid", clave)
			ctx.Sessionparams.Set("clientdata", clientdata.Data)
		}
	}
}

/* connect clientid and force cookie
func ForceUser(ctx *assets.Context, clientid int, longlogin int, origin string, source string, device string) {

	sitecontextname, _ := ctx.Sysparams.GetString("sitecontext")
	sitecontext := base.Sites.GetContext(sitecontextname)

	// NO podemos conectar un cliente que no existe
	clientdata, _ := graph.GetClient(sitecontext, clientid)
	if clientdata == nil {
		return
	}

	config := base.Sites.GetContext("kiwi-gr").Config

	cookiename, _ := config.GetString("cookiename")
	cookiedomain, _ := config.GetString("cookiedomain")

	sessionid := CreateSessionClient(ctx, clientid, longlogin, origin, source, device)
	ctx.Sessionparams.Set("clientid", clientid)
	ctx.Sessionparams.Set("clientdata", clientdata.Data)

	http.SetCookie(ctx.Writer, &http.Cookie{Name: cookiename, Value: sessionid, Path: "/", Domain: cookiedomain})

	//  fmt.Println("Forzando el cliente=", clientid, "con la cookie", cookiename, "=", sessionid)
}

/*
func CreateSessionUser(ctx *assets.Context, clientid int, longlogin int, origin string, source string, device string) string {

	config := base.Sites.GetContext("kiwi-gr").Config

	cookiesize, _ := config.GetInt("cookiesize")

	// usuario actualmente conectado (cargado al principio del hit por cookie), o nada
	clientcnx, _ := ctx.Sessionparams.GetInt("clientid")
	sessionid, _ := ctx.Sessionparams.GetString("sessionid")

	// Si ya estabamos conectados, crea una nueva sessionid
	if sessionid != "" && clientcnx != clientid {
		sessionid = ""
	}
	// crea sessionid
	ip, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)
	if sessionid == "" {
		// busca un ID disponible
		for {
			sessionid = util.CreateKey(cookiesize, -1)
			num, _ := client.KL_chefsesion.Count(xdominion.NewXCondition("clave", "=", sessionid))
			if num == 0 {
				break
			}
		}

		_, err := client.KL_chefsesion.Insert(xdominion.XRecord{
			"clave":       sessionid,
			"chef":        clientid,
			"fechainicio": time.Now(),
			"fechafin":    time.Now(),
			"ip":          ip,
			"longlogin":   longlogin,
			"origen":      origin + source,
			"device":      device,
		})
		if err != nil {
			fmt.Println("Error insertando sesion:", err)
		}
		ctx.Sessionparams.Set("sessionid", sessionid)
		ctx.Sessionparams.Set("clientid", clientid)
	} else {
		_, err := client.KL_chefsesion.Update(sessionid, xdominion.XRecord{
			"fechafin": time.Now(),
			"origen":   origin + source,
			"device":   device,
		})
		if err != nil {
			fmt.Println("Error modificando sesion:", err)
		}
	}

	ipdata := geoip.GetGeoData(ip)
	_, err := client.KL_chef.Update(clientid, xdominion.XRecord{
		"ultimopais":     ipdata.Country.IsoCode,
		"ultimaconexion": time.Now(),
		"intento":        0,
	})
	if err != nil {
		fmt.Println("Error modificando chef:", err)
	}
	return sessionid
}

func DestroySessionUser(ctx *assets.Context, sessionid string) {
	/*
	     $this->siteSesion = null;
	   $this->clientid = null;
	   $this->chefData = null;
	   SetCookie($this->base->config->cookiename, null, 0, '/', $this->base->config->COOKIEDOMAIN);

	   // comparte con Base
	   $this->base->clientid = $this->clientid;
	   $this->chefEntity->chefData = $this->chefData;
	   $this->base->siteSesion = $this->siteSesion;
	* /
}
*/
