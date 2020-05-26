// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package country

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID         = "country"
	VERSION          = "0.0.1"
	TRANSLATIONTHEME = "country"
)

func init() {
	m := &base.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Wiki", language.Spanish: "Wiki", language.French: "Wiki"},
		Needs:        []string{"base"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to ctx::database
func Setup(ds serverassets.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	buildTables(ctx)
	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	messages := []string{}

	ctx := ds.(*base.Datasource)
	// Needed modules: base and translation
	vc := base.ModuleInstalledVersion(ctx, "base")
	if vc == "" {
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/country.")
		return messages, nil
	}
	vc = base.ModuleInstalledVersion(ctx, "translation")
	if vc == "" {
		messages = append(messages, "xmodules/translation need to be installed before installing xmodules/country.")
		return messages, nil
	}

	translation.AddTheme(ctx, TRANSLATIONTHEME, "Countries", translation.SOURCETABLE, "", "name")

	// create tables
	messages = append(messages, createTables(ctx)...)

	// fill countries and translations
	messages = append(messages, loadTables(ctx, prefix)...)

	// Inserting into base-modules
	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(ctx, MODULEID, "List of official countries and ISO codes", VERSION)
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
