// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package wiki

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID = "wiki"
	VERSION  = "0.0.1"
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
// adds tables and caches to sitecontext::database
func Setup(ds applications.Datasource, prefix string) ([]string, error) {

	//	buildTables(ctx)
	//	createCache(ctx)
	ds.SetModule(MODULEID, VERSION)

	//	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	messages := []string{}

	// Needed modules: base and translation
	vc := base.ModuleInstalledVersion(ds, "base")
	if vc == "" {
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/wiki.")
		return messages, nil
	}

	// create tables
	messages = append(messages, createTables(ds)...)
	// fill data
	messages = append(messages, loadTables(ds)...)

	// Inserting into base-modules
	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(ds, MODULEID, "Wiki", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages, nil
}

func StartContext(ds applications.Datasource, ctx *context.Context) error {
	return nil
}
