// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package base

import (
	"fmt"

	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base/assets"
	"github.com/webability-go/xmodules/tools"
)

const (
	MODULEID = "base"
	VERSION  = "0.1.1"
)

var ModuleBase = assets.ModuleEntries{
	TryDatasource: TryDatasource,
}

func init() {
	messages = tools.BuildMessages(smessages)
	m := &Module{
		ID:      MODULEID,
		Version: VERSION,
		Languages: map[language.Tag]string{
			language.English: tools.Message(messages, "MODULENAME", language.English),
			language.Spanish: tools.Message(messages, "MODULENAME", language.Spanish),
			language.French:  tools.Message(messages, "MODULENAME", language.French),
		},
		Needs:         []string{},
		FSetup:        setup,
		FSynchronize:  synchronize,
		FStartContext: startContext,
	}
	ModulesList.Register(m)
}

// ======================================

// InitContext is called during the init phase to link the module with the system
// It must be called AFTER GetContainer
// adds tables and caches to sitecontext::database
// It should be called AFTER createContext
func setup(ds serverassets.Datasource, prefix string) ([]string, error) {

	linkTables(ds)
	ds.SetModule(MODULEID, VERSION)

	return []string{}, nil
}

func synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	result := []string{}
	tablename := "base_module"

	result = append(result, tools.Message(messages, "analyze", tablename))

	base_module := ds.GetTable(tablename)
	if base_module == nil {
		result = append(result, tools.Message(messages, "notable", tablename))
		return result, nil
	}
	num, err := base_module.Count(nil)
	if err != nil || num == 0 {
		if err != nil {
			result = append(result, tools.Message(messages, "tablenoexist", tablename, err))
		}
		err1 := base_module.Synchronize()
		if err1 != nil {
			result = append(result, tools.Message(messages, "tableerror", tablename, err1))
		} else {
			result = append(result, tools.Message(messages, "tablecreated", tablename))
		}
	} else {
		result = append(result, tools.Message(messages, "tablenotmodified", tablename))
	}

	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	// lets clone ds to begin a transaction
	cds := ds.CloneShell()
	_, err = cds.StartTransaction()
	if err != nil {
		result = append(result, err.Error())
		return result, err
	}

	err = AddModule(cds, MODULEID, tools.Message(messages, "MODULENAME"), VERSION)
	fmt.Println("Adds the module in table")
	if err == nil {
		result = append(result, tools.Message(messages, "modulemodified", MODULEID))
		result = append(result, tools.Message(messages, "commit"))
		cds.Commit()
		// TODO(Phil) should we also get the commit error if any?
	} else {
		result = append(result, tools.Message(messages, "rollback", err))
		cds.Rollback()
	}

	return result, nil
}

func startContext(ds serverassets.Datasource, ctx *serverassets.Context) error {
	return nil
}
