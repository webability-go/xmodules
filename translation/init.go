// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package translation

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation/assets"
)

const (
	MODULEID = "translation"
	VERSION  = "0.0.1"
)

var ModuleTranslation = assets.ModuleEntries{}

func init() {
	m := &base.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Translation tables", language.Spanish: "Tablas de traducción", language.French: "Tables de traduction"},
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
	ctx.SetModule(MODULEID, VERSION)

	return []string{}, nil
}

func Synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	messages := []string{}

	ctx := ds.(*base.Datasource)
	// Needed modules: base
	vc := base.ModuleInstalledVersion(ctx, "base")
	if vc == "" {
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/translation.")
		return messages, nil
	}
	vc1 := base.ModuleInstalledVersion(ctx, "user")
	vc2 := base.ModuleInstalledVersion(ctx, "userlink")
	if vc1 == "" && vc2 == "" {
		messages = append(messages, "xmodules/user or xmodules/userlink need to be installed before installing xmodules/translation.")
		return messages, nil
	}

	// create tables
	messages = append(messages, createTables(ctx)...)

	// Inserting into base-modules
	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(ctx, MODULEID, "Multilanguages translation tables for Xamboo", VERSION)
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
