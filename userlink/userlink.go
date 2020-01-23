// Package userlink contains the list of administrative user for the system, copied from a controller node
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package userlink

import (
	"fmt"
	"strconv"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "userlink"
	VERSION  = "1.0.2"
)

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func InitModule(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	sitecontext.Modules[MODULEID] = VERSION

	return nil
}

func SynchronizeModule(sitecontext *context.Context) []string {

	messages := []string{}

	// Needed modules: context and translation
	vc := context.ModuleInstalledVersion(sitecontext, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/userlink.")
		return messages
	}

	// create tables
	messages = append(messages, createTables(sitecontext)...)

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := context.AddModule(sitecontext, MODULEID, "Administration users link", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages
}

func SynchroUsers(sitecontext *context.Context, fromcontext *context.Context) []string {

	msg := []string{}

	// load from origin context
	table, ok := fromcontext.Tables["user_user"]
	if !ok {
		return []string{"Error: the origin table user_user does not exist."}
	}
	users, _ := table.SelectAll(nil, xdominion.XFieldSet{"key", "status", "name", "mail"})
	keys := map[int]bool{}
	if users != nil {
		for _, u := range *users {
			key, _ := u.GetInt("key")
			_, err := sitecontext.Tables["user_user"].Upsert(key, u)
			if err != nil {
				msg = append(msg, "Error adding user:"+fmt.Sprint(err))
			}
			keys[key] = true
		}
	}
	// extra users to delete
	localusers, _ := table.SelectAll(nil, xdominion.XFieldSet{"key"})
	localkeys := map[int]bool{}
	if localusers != nil {
		for _, u := range *localusers {
			key, _ := u.GetInt("key")
			localkeys[key] = true
		}
	}
	// diff localusers - users
	for k := range keys {
		if localkeys[k] {
			delete(localkeys, k)
		}
	}
	// set status to "deleted" to X (not available anymore), they may be used by the local code
	for k := range localkeys {
		sitecontext.Tables["user_user"].Update(k, xdominion.XRecord{"status": "X"})
	}

	cnt, _ := sitecontext.Tables["user_user"].Count(nil)
	msg = append(msg, strconv.Itoa(cnt)+" admin users synchronized")
	return msg
}
