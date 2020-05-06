// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package user

import (
	"strings"

	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "user"
	VERSION  = "0.0.1"
)

func init() {
	m := &context.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Users for administration", language.Spanish: "Usuarios para administración", language.French: "Utilisateurs pour l'administration"},
		Needs:        []string{"context"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	context.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(sitecontext *context.Context, databasename string) (string, error) {

	buildTables(sitecontext, databasename)
	createCache(sitecontext)
	sitecontext.SetModule(MODULEID, VERSION)

	go buildCache(sitecontext)

	return "", nil
}

func Synchronize(sitecontext *context.Context, databasename string) (string, error) {

	messages := []string{}

	// Needed modules: context and translation
	vc := context.ModuleInstalledVersion(sitecontext, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/user.")
		return strings.Join(messages, "\n"), nil
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

	return strings.Join(messages, "\n"), nil
}
