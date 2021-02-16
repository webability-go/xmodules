// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package usda

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation"
	"github.com/webability-go/xmodules/usda/assets"
)

const (
	MODULEID         = "usda"
	VERSION          = "0.0.1"
	TRANSLATIONTHEME = "nutrient"
)

var ModuleUsda = assets.ModuleEntries{}

func init() {
	m := &base.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "USDA Database", language.Spanish: "Base de datos de la USDA", language.French: "Base de données de la USDA"},
		Needs:        []string{"base"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(ds applications.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	buildTables(ctx)
	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	translation.AddTheme(ctx, TRANSLATIONTHEME, "USDA nutrients", translation.SOURCETABLE, "", "name,tag")

	messages := []string{}
	messages = append(messages, createTables(ctx)...)

	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(ctx, MODULEID, "List of USDA food and nutrients", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	messages = append(messages, loadTables(ctx, prefix)...)
	messages = append(messages, buildCache(ctx)...)

	return messages, nil
}

func StartContext(ds applications.Datasource, ctx *context.Context) error {
	return nil
}
