package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
)

func GetSession(ds *base.Datasource, sessionid string) *xdominion.XRecord {

	user_session := ds.GetTable("user_session")
	if user_session == nil {
		ds.Log("xmodules::user::GetSession: Error, the user_session table is not available on this datasource")
		return nil
	}

	data, _ := user_session.SelectOne(sessionid)
	return data
}

func CreateSession(ds *base.Datasource, key int, sessionid string, IP string, origin string, device string) string {

	user_session := ds.GetTable("user_session")
	if user_session == nil {
		ds.Log("xmodules::user::CreateSession: Error, the user_session table is not available on this datasource")
		return ""
	}

	// Lets see if we can reuse the sessionid
	if sessionid != "" {
		// load session to see if it fits, or not
		sessiondata := GetSession(ds, sessionid)
		if sessiondata != nil {
			userkey, _ := sessiondata.GetInt("user")
			if userkey == key {
				// YES, we can apply to this this session
				_, err := user_session.Update(sessionid, xdominion.XRecord{
					"lastmodif": time.Now(),
					"ip":        IP,
					"origen":    origin,
					"device":    device,
				})
				if err == nil {
					return sessionid
				}
				fmt.Println("Error inserting session:", err)
				sessionid = ""
			}
		}
	}

	// busca un ID disponible
	for {
		sessionid = tools.UUID()
		num, _ := user_session.Count(xdominion.NewXCondition("key", "=", sessionid))
		if num == 0 {
			break
		}
	}

	_, err := user_session.Insert(xdominion.XRecord{
		"key":          sessionid,
		"user":         key,
		"creationdate": time.Now(),
		"lastmodif":    time.Now(),
		"ip":           IP,
		"longlogin":    1,
		"origin":       origin,
		"device":       device,
	})
	if err != nil {
		fmt.Println("Error inserting sesion:", err)
		return ""
	}
	return sessionid
}

func CloseSession(ds *base.Datasource, sessionid string) string {

	user_session := ds.GetTable("user_session")
	if user_session == nil {
		ds.Log("xmodules::user::CreateSession: Error, the user_session table is not available on this datasource")
		return ""
	}
	// invaluda sessionid
	newsessionid := strings.Replace(sessionid, "-", "=", -1)
	_, err := user_session.Update(sessionid, xdominion.XRecord{
		"lastmodif": time.Now(),
		"key":       newsessionid,
	})
	if err != nil {
		fmt.Println("Error closing session:", err)
		return sessionid
	}
	return ""
}
