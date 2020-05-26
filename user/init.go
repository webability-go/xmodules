// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package user

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/user/assets"
)

const (
	MODULEID = "user"
	VERSION  = "0.0.1"
)

var ModuleUser = assets.ModuleEntries{
	VerifyUserSession: VerifyUserSession,
	GetUsersList:      GetUsersList,
}

func init() {
	m := &base.Module{
		ID:            MODULEID,
		Version:       VERSION,
		Languages:     map[language.Tag]string{language.English: "Users for administration", language.Spanish: "Usuarios para administración", language.French: "Utilisateurs pour l'administration"},
		Needs:         []string{"base"},
		FSetup:        Setup,        // Called once at the main system startup, once PER CREATED xmodule CONTEXT (if set)
		FSynchronize:  Synchronize,  // Called only to create/rebuild database objects and others on demand (if set)
		FStartContext: StartContext, // Called each time a new Server context is created  (if set)
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(ds serverassets.Datasource, prefix string) ([]string, error) {

	lds := ds.(*base.Datasource)
	buildTables(lds)
	createCache(lds)
	lds.SetModule(MODULEID, VERSION)

	go buildCache(lds)

	return []string{}, nil
}

func Synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	messages := []string{}

	lds := ds.(*base.Datasource)
	// Needed modules: context and translation
	vc := base.ModuleInstalledVersion(ds, "base")
	if vc == "" {
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/user.")
		return messages, nil
	}

	// create tables
	messages = append(messages, createTables(lds)...)
	// fill super admin
	messages = append(messages, loadTables(lds)...)

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := base.AddModule(ds, MODULEID, "Administration users", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages, nil
}

func StartContext(ds serverassets.Datasource, ctx *serverassets.Context) error {

	// TODO(phil) implement site, device (from CTX)

	lds := ds.(*base.Datasource)
	VerifyUserSession(ctx, lds, "xmodules", "pc")
	return nil
}
