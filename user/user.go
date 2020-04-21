// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package user

import (
	"fmt"
	"strconv"
	"time"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "user"
	VERSION  = "2.0.0"
)

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func InitModule(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	createCache(sitecontext)
	sitecontext.SetModule(MODULEID, VERSION)

	go buildCache(sitecontext)

	return nil
}

func SynchronizeModule(sitecontext *context.Context) []string {

	messages := []string{}

	// Needed modules: context and translation
	vc := context.ModuleInstalledVersion(sitecontext, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/user.")
		return messages
	}

	// create tables
	messages = append(messages, createTables(sitecontext)...)
	// fill super admin
	messages = append(messages, loadTables(sitecontext)...)

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := context.AddModule(sitecontext, MODULEID, "Administration users", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetUser(sitecontext *context.Context, key int) *StructureUser {

	cache := sitecontext.GetCache("user:users")
	if cache == nil {
		sitecontext.Log("main", "xmodules::user::GetUser: Error, the user cache is not available on this site context")
		return nil
	}

	data, _ := cache.Get(strconv.Itoa(key))
	if data == nil {
		sm := CreateStructureUserByKey(sitecontext, key)
		if sm == nil {
			sitecontext.Log("graph", "xmodules::user::GetUser: There is no user created: ", key)
			return nil
		}
		cache.Set(strconv.Itoa(key), sm)
		return sm.(*StructureUser)
	}
	return data.(*StructureUser)
}

func GetUserByCredentials(sitecontext *context.Context, username string, password string) *StructureUser {

	user_user := sitecontext.GetTable("user_user")
	if user_user == nil {
		sitecontext.Log("xmodules::user::GetUserByCredentials: Error, the user_user table is not available on this context")
		return nil
	}
	fmt.Println("GetUserByCredentials:", username, password)
	data, _ := user_user.SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("un", "=", username),
		xdominion.NewXCondition("pw", "=", password, "and"),
		xdominion.NewXCondition("status", "!=", "X", "and"),
	})
	if data == nil {
		return nil
	}
	fmt.Println("GetUserByCredentials:", data)

	sm := CreateStructureUserByData(sitecontext, data)
	if sm == nil {
		sitecontext.Log("graph", "xmodules::user::GetUser: There is no user created: ", fmt.Sprint(data))
		return nil
	}
	return sm.(*StructureUser)
}

func GetSession(sitecontext *context.Context, sessionid string) *xdominion.XRecord {

	user_session := sitecontext.GetTable("user_session")
	if user_session == nil {
		sitecontext.Log("xmodules::user::GetSession: Error, the user_session table is not available on this context")
		return nil
	}

	data, _ := user_session.SelectOne(sessionid)
	return data
}

func CreateSession(sitecontext *context.Context, keysize int, key int, sessionid string, IP string, origin string, device string) string {

	user_session := sitecontext.GetTable("user_session")
	if user_session == nil {
		sitecontext.Log("xmodules::user::CreateSession: Error, the user_session table is not available on this context")
		return ""
	}
	fmt.Println("CreateSession", keysize, key, sessionid, IP, origin, device)
	// Lets see if we can reuse the sessionid
	if sessionid != "" {
		// load session to see if it fits, or not
		sessiondata := GetSession(sitecontext, sessionid)
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
				fmt.Println("Error insertint session:", err)
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

func CloseSession(sitecontext *context.Context, sessionid string) string {

	user_session := sitecontext.GetTable("user_session")
	if user_session == nil {
		sitecontext.Log("xmodules::user::CreateSession: Error, the user_session table is not available on this context")
		return ""
	}
	// invaluda sessionid
	sessionid[len(sessionid)-1] = 'Z'
	_, err := user_session.Update(sessionid, xdominion.XRecord{
		"lastmodif": time.Now(),
		"sessionid": sessionid,
	})
	if err == nil {
		return sessionid
	}

}
