// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package useradmin

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/useradmin/assets"
)

const (
	MODULEID = "useradmin"
	VERSION  = "0.0.1"
)

var Needs = []string{"base", "user", "adminmenu"}
var ModuleUserAdmin = assets.ModuleEntries{}

func init() {
	messages = tools.BuildMessages(smessages)
	m := &base.Module{
		ID:      MODULEID,
		Version: VERSION,
		Languages: map[language.Tag]string{
			language.English: tools.Message(messages, "MODULENAME", language.English),
			language.Spanish: tools.Message(messages, "MODULENAME", language.Spanish),
			language.French:  tools.Message(messages, "MODULENAME", language.French),
		},
		Needs:         Needs,
		FSetup:        Setup,        // Called once at the main system startup, once PER CREATED xmodule CONTEXT (if set)
		FSynchronize:  Synchronize,  // Called only to create/rebuild database objects and others on demand (if set)
		FStartContext: StartContext, // Called each time a new Server context is created  (if set)
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(ds serverassets.Datasource, prefix string) ([]string, error) {

	// no tables on this module
	ds.SetModule(MODULEID, VERSION)

	return []string{}, nil
}

func Synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	result := []string{}

	ok, res := base.VerifyNeeds(ds, Needs)
	result = append(result, res...)
	if !ok {
		return result, nil
	}

	installed := base.ModuleInstalledVersion(ds, MODULEID)

	cds := ds.CloneShell()
	_, err := cds.StartTransaction()
	if err != nil {
		result = append(result, err.Error())
		return result, err
	}

	var r []string
	// installation or upgrade ?
	if installed != "" {
		err, r = upgrade(cds, installed)
	} else {
		err, r = install(cds)
	}
	result = append(result, r...)
	if err == nil {
		err = base.AddModule(cds, MODULEID, tools.Message(messages, "MODULENAME"), VERSION)
		if err == nil {
			result = append(result, tools.Message(messages, "modulemodified", MODULEID))
			result = append(result, tools.Message(messages, "commit"))
			err = cds.Commit()
			if err != nil {
				result = append(result, err.Error())
			}
		} else {
			result = append(result, tools.Message(messages, "rollback", err))
			err = cds.Rollback()
			if err != nil {
				result = append(result, err.Error())
			}
		}
	}

	return result, nil
}

func StartContext(ds serverassets.Datasource, ctx *serverassets.Context) error {
	return nil
}
