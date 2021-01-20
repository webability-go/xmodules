// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package adminmenu

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/adminmenu/assets"
	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID   = "adminmenu"
	VERSION    = "0.0.1"
	DATASOURCE = "adminmenudatasource"
)

var ModuleAdminMenu = assets.ModuleEntries{
	AddGroup:  AddGroup,
	GetGroup:  GetGroup,
	AddOption: AddOption,
	GetOption: GetOption,
	GetMenu:   GetMenu,
}

func init() {
	m := &base.Module{
		ID:            MODULEID,
		Version:       VERSION,
		Languages:     map[language.Tag]string{language.English: "Administration menu", language.Spanish: "Menu de administración", language.French: "Menu pour l'administration"},
		Needs:         []string{"base", "user"},
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
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/adminmenu.")
		return messages, nil
	}

	vc = base.ModuleInstalledVersion(ds, "user")
	if vc == "" {
		messages = append(messages, "xmodules/user need to be installed before installing xmodules/adminmenu.")
		return messages, nil
	}

	// create tables
	messages = append(messages, createTables(lds)...)
	// fill super admin
	messages = append(messages, loadTables(lds)...)

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := base.AddModule(ds, MODULEID, "Administration Admin menu", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages, nil
}

func StartContext(ds serverassets.Datasource, ctx *serverassets.Context) error {

	return nil
}
