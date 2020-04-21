package user

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/webability-go/xcore/v2"
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
		table := moduletables[tbl]()
		table.SetBase(sitecontext.GetDatabase(databasename))
		sitecontext.SetTable(tbl, table)
	}
}

func createCache(sitecontext *context.Context) []string {

	sitecontext.SetCache("user:users", xcore.NewXCache("user:users", 0, 0))

	return []string{}
}

func buildCache(sitecontext *context.Context) []string {

	user_user := sitecontext.GetTable("user_user")
	if user_user == nil {
		return []string{"xmodules::user::buildCache: Error, the user_user table is not available on this context"}
	}
	cache := sitecontext.GetCache("user:users")
	if cache == nil {
		return []string{"xmodules::user::buildCache: Error, the user cache is not available on this site context"}
	}

	// Loads all data in XCache
	users, _ := user_user.SelectAll()

	if users != nil {
		for _, m := range *users {
			// creates structure on language
			str := CreateStructureUserByData(sitecontext, m.Clone())
			key, _ := m.GetString("key")
			cache.Set(key, str)
		}
	}

	return []string{}
}

func createTables(sitecontext *context.Context) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {

		table := sitecontext.GetTable(tbl)
		if table == nil {
			return []string{"xmodules::user::createTables: Error, the table is not available on this context:" + tbl}
		}

		messages = append(messages, "Analysing "+tbl+" table.")
		num, err := table.Count(nil)
		if err != nil || num == 0 {
			err1 := table.Synchronize()
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

	user_user := sitecontext.GetTable("user_user")
	if user_user == nil {
		return []string{"xmodules::user::createTables: Error, the table user_user is not available on this context"}
	}

	// insert super admin
	_, err := user_user.Upsert(1, xdominion.XRecord{
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
		sitecontext.Log("main", "Error inserting admin user", err)
		return []string{"xmodules::user::loadTables: Error upserting the admin user"}
	}
	return []string{
		fmt.Sprint(user_user.Count(nil)) + " admin user added",
	}
}
