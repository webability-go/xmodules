// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package user

import (
	"strconv"

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
