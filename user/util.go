package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/base"
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

func linkTables(ds applications.Datasource) {

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(ds.GetDatabase())
		ds.SetTable(tbl, table)
	}
}

func createCache(ds applications.Datasource) []string {

	ds.SetCache("user:users", xcore.NewXCache("user:users", 0, 0))

	return []string{}
}

func buildCache(ds applications.Datasource) []string {

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		return []string{"xmodules::user::buildCache: Error, the user_user table is not available on this datasource"}
	}
	cache := ds.GetCache("user:users")
	if cache == nil {
		return []string{"xmodules::user::buildCache: Error, the user cache is not available on this site datasource"}
	}

	// Loads all data in XCache
	users, _ := user_user.SelectAll()

	if users != nil {
		for _, m := range *users {
			// creates structure on language
			str := CreateStructureUserByData(ds, m.Clone())
			key, _ := m.GetString("key")
			cache.Set(key, str)
		}
	}

	return []string{}
}

func synchroTables(ds applications.Datasource, oldversion string) (error, []string) {

	result := []string{}

	for _, tbl := range moduletablesorder {

		err, res := base.SynchroTable(ds, tbl)
		result = append(result, res...)
		if err != nil {
			return err, result
		}
	}

	if oldversion < "0.0.1" {
		// do things
	}

	return nil, result
}

func install(ds applications.Datasource) (error, []string) {

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		return errors.New("error"), []string{"xmodules::user::createTables: Error, the table user_user is not available on this datasource"}
	}

	bmd5 := md5.Sum([]byte("manager"))
	// insert super admin
	_, err := user_user.Upsert(1, xdominion.XRecord{
		"key":          1,
		"status":       "S",
		"name":         "System Manager",
		"un":           "system",
		"pw":           hex.EncodeToString(bmd5[:]),
		"mail":         "hostmaster@yoursite.com",
		"sex":          "M",
		"creationdate": time.Now(),
		"lastmodif":    time.Now(),
	})
	if err != nil {
		ds.Log("main", "Error inserting admin user", err)
		return errors.New("Error"), []string{"xmodules::user::loadTables: Error upserting the admin user"}
	}
	return nil, []string{
		fmt.Sprint(user_user.Count(nil)) + " admin user added",
	}
}

func upgrade(ds applications.Datasource, oldversion string) (error, []string) {

	if oldversion < "0.0.1" {
		// do things
	}

	return nil, []string{}
}
