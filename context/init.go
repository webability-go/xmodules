// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package context

import (
	"strings"

	"golang.org/x/text/language"
)

const (
	MODULEID = "context"
	VERSION  = "0.0.1"
)

func init() {
	m := &Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Contexts", language.Spanish: "Contextos", language.French: "Contextes"},
		Needs:        []string{},
		FSetup:       setupModule,
		FSynchronize: synchronizeModule,
	}
	ModulesList.Register(m)
}

// ======================================

// InitContext is called during the init phase to link the module with the system
// It must be called AFTER GetContainer
// adds tables and caches to sitecontext::database
// It should be called AFTER createContext
func setupModule(context *Context, db string) (string, error) {

	buildTables(context, db)
	buildCache(context)
	context.SetModule(MODULEID, VERSION)

	return "", nil
}

func synchronizeModule(context *Context, db string) (string, error) {

	messages := []string{}
	messages = append(messages, "Analysing context_module table.")

	context_module := context.GetTable("context_module")
	if context_module == nil {
		messages = append(messages, "Critical Error: the context table context_module does not exist !!!: ")
		return strings.Join(messages, "\n"), nil
	}
	num, err := context_module.Count(nil)
	if err != nil || num == 0 {
		err1 := context_module.Synchronize()
		if err1 != nil {
			messages = append(messages, "The table context_module was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table context_module was created (again)")
		}
	} else {
		messages = append(messages, "The table context_module was not created because it contains data.")
	}

	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err = AddModule(context, MODULEID, "Contexts and Modules for Xamboo", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the context_module table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the context_module table: "+err.Error())
	}
	return strings.Join(messages, "\n"), nil
}
