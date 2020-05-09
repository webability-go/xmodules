// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package country

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID         = "country"
	VERSION          = "0.0.1"
	TRANSLATIONTHEME = "country"
)

func init() {
	m := &context.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Wiki", language.Spanish: "Wiki", language.French: "Wiki"},
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
	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ctx *context.Context, filespath string) ([]string, error) {

	messages := []string{}

	// Needed modules: context and translation
	vc := context.ModuleInstalledVersion(ctx, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/country.")
		return messages, nil
	}
	vc = context.ModuleInstalledVersion(ctx, "translation")
	if vc == "" {
		messages = append(messages, "xmodules/translation need to be installed before installing xmodules/country.")
		return messages, nil
	}

	translation.AddTheme(ctx, TRANSLATIONTHEME, "Countries", translation.SOURCETABLE, "", "name")

	// create tables
	messages = append(messages, createTables(ctx)...)

	// fill countries and translations
	messages = append(messages, loadTables(ctx, filespath)...)

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := context.AddModule(ctx, MODULEID, "List of official countries and ISO codes", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages, nil
}
