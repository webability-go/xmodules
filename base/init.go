// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package base

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/assets"
)

const (
	MODULEID = "base"
	VERSION  = "0.1.0"
)

func init() {
	m := &Module{
		ID:            MODULEID,
		Version:       VERSION,
		Languages:     map[language.Tag]string{language.English: "XModules base", language.Spanish: "Base XModules", language.French: "Base XModules"},
		Needs:         []string{},
		FSetup:        setupModule,
		FSynchronize:  synchronizeModule,
		FStartContext: startcontext,
	}
	ModulesList.Register(m)
}

// ======================================

// InitContext is called during the init phase to link the module with the system
// It must be called AFTER GetContainer
// adds tables and caches to sitecontext::database
// It should be called AFTER createContext
func setupModule(ds assets.Datasource, prefix string) ([]string, error) {

	buildTables(ds)
	buildCache(ds)
	ds.SetModule(MODULEID, VERSION)

	return []string{}, nil
}

func synchronizeModule(ds assets.Datasource, prefix string) ([]string, error) {

	messages := []string{}
	messages = append(messages, "Analysing base_module table.")

	base_module := ds.GetTable("base_module")
	if base_module == nil {
		messages = append(messages, "Critical Error: the base table base_module does not exist !!!: ")
		return messages, nil
	}
	num, err := base_module.Count(nil)
	if err != nil || num == 0 {
		err1 := base_module.Synchronize()
		if err1 != nil {
			messages = append(messages, "The table base_module was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table base_module was created (again)")
		}
	} else {
		messages = append(messages, "The table context_module was not created because it contains data.")
	}

	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err = AddModule(ds, MODULEID, "Contexts and Modules for Xamboo", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the base_module table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the base_module table: "+err.Error())
	}
	return messages, nil
}

func startcontext(ds assets.Datasource, ctx *assets.Context) error {
	return nil
}
