package user

import (
	"fmt"
	"time"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{
	"user_user",
	"user_accessgroup", "user_access", "user_accessextended",
	"user_profile", "user_profileaccess", "user_profileaccessextended",
	"user_useraccess", "user_useraccessextended",
	"user_userprofile",
	"user_parameter",
	"user_session", "user_sessionhistory",
}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{
	"user_user":                  userUser,
	"user_accessgroup":           userAccessGroup,
	"user_access":                userAccess,
	"user_accessextended":        userAccessExtended,
	"user_profile":               userProfile,
	"user_profileaccess":         userProfileAccess,
	"user_profileaccessextended": userProfileAccessExtended,
	"user_useraccess":            userUserAccess,
	"user_useraccessextended":    userUserAccessExtended,
	"user_userprofile":           userUserProfile,
	"user_parameter":             userParameter,
	"user_session":               userSession,
	"user_sessionhistory":        userSessionHistory,
}

func buildTables(sitecontext *context.Context, databasename string) {

	for _, tbl := range moduletablesorder {
		sitecontext.Tables[tbl] = moduletables[tbl]()
		sitecontext.Tables[tbl].SetBase(sitecontext.Databases[databasename])
	}
}

func buildCache(sitecontext *context.Context) []string {

	// Loads all data in XCache
	users, _ := sitecontext.Tables["user_user"].SelectAll()

	sitecontext.Caches["user:users"] = xcore.NewXCache("user:users", 0, 0)

	if users != nil {
		for _, m := range *users {
			// creates structure on language
			str := CreateStructureUserByData(sitecontext, m.Clone())
			key, _ := m.GetString("key")
			sitecontext.Caches["user:users"].Set(key, str)
		}
	}

	return []string{}
}

func createTables(sitecontext *context.Context) []string {

	messages := []string{}

	for tbl := range moduletables {
		messages = append(messages, "Analysing "+tbl+" table.")
		num, err := sitecontext.Tables[tbl].Count(nil)
		if err != nil || num == 0 {
			err1 := sitecontext.Tables[tbl].Synchronize()
			if err1 != nil {
				messages = append(messages, "The table "+tbl+" was not created: "+err1.Error())
			} else {
				messages = append(messages, "The table "+tbl+" was created (again)")
			}
		} else {
			messages = append(messages, "The table "+tbl+" was not created because it contains data.")
		}
	}

	return messages
}

func loadTables(sitecontext *context.Context) []string {

	// insert super admin
	_, err := sitecontext.Tables["user_user"].Upsert(1, xdominion.XRecord{
		"key":          1,
		"status":       "A",
		"name":         "System Manager",
		"un":           "system",
		"pw":           "manager",
		"mail":         "hostmaster@yoursite.com",
		"sex":          "M",
		"creationdate": time.Now(),
		"lastmodif":    time.Now(),
	})
	if err != nil {

	}
	return []string{
		fmt.Sprint(sitecontext.Tables["user_user"].Count(nil)) + " admin user added",
	}
}
