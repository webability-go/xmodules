// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package suggestions

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID = "suggestions"
	VERSION  = "0.0.1"
)

func init() {
	m := &base.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Suggestion tables", language.Spanish: "Tablas de sugerencias", language.French: "Tables de suggestions"},
		Needs:        []string{"base"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(ds serverassets.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	buildTables(ctx)
	//	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	//	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	messages := []string{}

	ctx := ds.(*base.Datasource)
	// Needed modules: base and translation
	vc := base.ModuleInstalledVersion(ctx, "base")
	if vc == "" {
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/suggestions.")
		return messages, nil
	}

	// create tables
	messages = append(messages, createTables(ctx)...)
	// fill super admin
	messages = append(messages, loadTables(ctx)...)

	// Inserting into base-modules
	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(ctx, MODULEID, "Suggestions", VERSION)
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
