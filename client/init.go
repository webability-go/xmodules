// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package client

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID = "client"
	VERSION  = "0.0.1"
)

func init() {
	m := &base.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Clients", language.Spanish: "Clientes", language.French: "Clients"},
		Needs:        []string{"base"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(ds serverassets.Datasource, prefix string) ([]string, error) {

	dts := ds.(*base.Datasource)
	buildTables(dts)
	createCache(dts)
	dts.SetModule(MODULEID, VERSION)

	go buildCache(dts)

	return []string{}, nil
}

func Synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	messages := []string{}

	dts := ds.(*base.Datasource)
	// Needed modules: base and translation
	vc := base.ModuleInstalledVersion(dts, "base")
	if vc == "" {
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/client.")
		return messages, nil
	}

	// create tables
	messages = append(messages, createTables(dts)...)
	// fill super admin
	messages = append(messages, loadTables(dts)...)

	// Inserting into base-modules
	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(dts, MODULEID, "Clients", VERSION)
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
