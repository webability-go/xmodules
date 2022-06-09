// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package base

import (
	"embed"
	"sync"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/tools"
)

const (
	MODULEID = "base"
	VERSION  = "0.1.1"
)

var Needs = []string{}

// var Containers = &ContainersList{}
var Containers *ContainersList
var ContainersLock sync.RWMutex

//go:embed languages/*.language
var fsmessages embed.FS
var messages *map[language.Tag]*xcore.XLanguage

func init() {
	// Creates the Containers into the init() to avoid the GO 1.16.1 error on creating global variables on plugins and not using them directly
	Containers = &ContainersList{}
	messages = tools.BuildMessagesFS(fsmessages, "languages")
	m := &Module{
		ID:      MODULEID,
		Version: VERSION,
		Languages: map[language.Tag]string{
			language.English: tools.Message(messages, "MODULENAME", language.English),
			language.Spanish: tools.Message(messages, "MODULENAME", language.Spanish),
			language.French:  tools.Message(messages, "MODULENAME", language.French),
		},
		Needs:         Needs,
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
func setup(ds applications.Datasource, prefix string) ([]string, error) {

	linkTables(ds)
	ds.SetModule(MODULEID, VERSION)

	return []string{}, nil
}

func synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	result := []string{}

	ok, res := VerifyNeeds(ds, Needs)
	result = append(result, res...)
	if !ok {
		return result, nil
	}

	installed := ModuleInstalledVersion(ds, MODULEID)

	// synchro tables
	err, r := synchroTables(ds, installed)
	result = append(result, r...)
	if err != nil {
		return result, err
	}

	// The rest of the process with a transaction
	// lets clone ds to begin a transaction
	cds := ds.CloneShell()
	_, err = cds.StartTransaction()
	if err != nil {
		result = append(result, err.Error())
		return result, err
	}

	// installation or upgrade ?
	if installed != "" {
		err, r = upgrade(cds, installed)
	} else {
		err, r = install(cds)
	}
	result = append(result, r...)
	if err == nil {
		err = AddModule(cds, MODULEID, tools.Message(messages, "MODULENAME"), VERSION)
		if err == nil {
			result = append(result, tools.Message(messages, "modulemodified", MODULEID))
			result = append(result, tools.Message(messages, "commit"))
			err = cds.Commit()
			if err != nil {
				result = append(result, err.Error())
			}
			return result, err
		}
	}
	result = append(result, tools.Message(messages, "rollback", err))
	err1 := cds.Rollback()
	if err1 != nil {
		result = append(result, err1.Error())
	}

	return result, err
}

func startContext(ds applications.Datasource, ctx *context.Context) error {
	return nil
}
