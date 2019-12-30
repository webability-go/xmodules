// Package userlink contains the list of administrative user for the system, copied from a controller node
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package userlink

import (
	"fmt"

	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "userlink"
	VERSION  = "1.0.0"
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

	// load from origin context

	// upsert all users in this context
	return []string{
		fmt.Sprint(sitecontext.Tables["user_user"].Count(nil)) + " admin users synchronized",
	}
}
