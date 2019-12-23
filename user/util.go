package user

import (
	"github.com/webability-go/xcore"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

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

	for tbl, fct := range moduletables {
		sitecontext.Tables[tbl] = fct()
		sitecontext.Tables[tbl].SetBase(sitecontext.Databases[databasename])
	}
}

func buildCache(sitecontext *context.Context) []string {

	// Loads all data in XCache
	users, _ := sitecontext.Tables["user_user"].SelectAll()

	sitecontext.Caches["user:users"] = xcore.NewXCache("user:users", 0, 0)

	for _, m := range *users {
		// creates structure on language
		str := CreateStructureUserByData(sitecontext, m.Clone())
		key, _ := m.GetString("key")
		sitecontext.Caches["user:users"].Set(key, str)
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

	return []string{}
}
