// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package usda

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID         = "usda"
	VERSION          = "0.0.1"
	TRANSLATIONTHEME = "nutrient"
)

func init() {
	m := &context.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "USDA Database", language.Spanish: "Base de datos de la USDA", language.French: "Base de données de la USDA"},
		Needs:        []string{"context"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	context.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(ctx *context.Context, prefix string) ([]string, error) {

	buildTables(ctx)
	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(sitecontext *context.Context, filespath string) ([]string, error) {

	translation.AddTheme(sitecontext, TRANSLATIONTHEME, "USDA nutrients", translation.SOURCETABLE, "", "name,tag")

	messages := []string{}
	messages = append(messages, createTables(sitecontext)...)

	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := context.AddModule(sitecontext, MODULEID, "List of USDA food and nutrients", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	messages = append(messages, loadTables(sitecontext, filespath)...)
	messages = append(messages, buildCache(sitecontext)...)

	return messages, nil
}
