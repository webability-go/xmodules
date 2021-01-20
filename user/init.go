// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package user

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user/assets"
)

const (
	MODULEID = "user"
	VERSION  = "0.0.1"
)

var Needs = []string{"base"}
var ModuleUser = assets.ModuleEntries{
	GetAccessesCount: GetCountAccesses,
	GetAccessesList:  GetAccessesList,
}

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

	buildTables(ds)
	createCache(ds)
	ds.SetModule(MODULEID, VERSION)

	go buildCache(ds)

	return []string{}, nil
}

func Synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	result := []string{}

	// Needed modules: context and translation
	vc := base.ModuleInstalledVersion(ds, "base")
	if vc == "" {
		result = append(result, "xmodules/base need to be installed before installing xmodules/user.")
		return result, nil
	}

	// create tables
	r, err := createTables(ds)
	result = append(result, r...)

	cds := ds.CloneShell()
	_, err = cds.StartTransaction()
	if err != nil {
		result = append(result, err.Error())
		return result, err
	}

	// fill super admin
	r, err = loadTables(cds)
	result = append(result, r...)

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err = base.AddModule(cds, MODULEID, "Administration users", VERSION)
	if err == nil {
		result = append(result, "The entry "+MODULEID+" was modified successfully in the modules table.")
		result = append(result, tools.Message(messages, "commit"))
		cds.Commit()
	} else {
		result = append(result, tools.Message(messages, "rollback", err))
		cds.Rollback()
	}

	return result, nil
}

func StartContext(ds serverassets.Datasource, ctx *serverassets.Context) error {

	// TODO(phil) implement site, device (from CTX)

	lds := ds.(*base.Datasource)
	VerifyUserSession(ctx, lds, "xmodules", "pc")
	return nil
}
