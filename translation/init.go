// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package translation

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "translation"
	VERSION  = "0.0.1"
)

func init() {
	m := &context.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Translation tables", language.Spanish: "Tablas de traducción", language.French: "Tables de traduction"},
		Needs:        []string{"context"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	context.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to ctx::database
func Setup(ctx *context.Context, prefix string) ([]string, error) {

	buildTables(ctx)
	ctx.SetModule(MODULEID, VERSION)

	return []string{}, nil
}

func Synchronize(ctx *context.Context, prefix string) ([]string, error) {

	messages := []string{}

	// Needed modules: context
	vc := context.ModuleInstalledVersion(ctx, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/translation.")
		return messages, nil
	}
	vc1 := context.ModuleInstalledVersion(ctx, "user")
	vc2 := context.ModuleInstalledVersion(ctx, "userlink")
	if vc1 == "" && vc2 == "" {
		messages = append(messages, "xmodules/user or xmodules/userlink need to be installed before installing xmodules/translation.")
		return messages, nil
	}

	// create tables
	messages = append(messages, createTables(ctx)...)

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := context.AddModule(ctx, MODULEID, "Multilanguages translation tables for Xamboo", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages, nil
}
