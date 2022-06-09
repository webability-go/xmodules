// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package stat

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID = "stat"
	VERSION  = "0.0.1"
)

func init() {
	m := &base.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Stat logs", language.Spanish: "Bitacoras", language.French: "Registres"},
		Needs:        []string{"base"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to ctx::database
func Setup(ds applications.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	buildTables(ctx, prefix)
	//	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	//	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	messages := []string{}

	ctx := ds.(*base.Datasource)
	// Needed modules: base and translation
	vc := base.ModuleInstalledVersion(ctx, "base")
	if vc == "" {
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/stat.")
		return messages, nil
	}

	for i := 1; i < 13; i++ {
		m := ""
		if i < 10 {
			m = "0"
		}
		m += strconv.Itoa(i)

		messages = append(messages, "Analysing "+prefix+"stat_"+m+" table.")
		table := ctx.GetTable(prefix + "stat_" + m)
		if table == nil {
			messages = append(messages, "xmodules::stat::SynchronizeModule: Error, the table does not exist in the base: "+prefix+"stat_"+m)
			return messages, nil
		}

		num, err := table.Count(nil)
		if err != nil || num == 0 {
			err1 := table.Synchronize()
			if err1 != nil {
				messages = append(messages, "The table "+prefix+"stat_"+m+" was not created: "+err1.Error())
			} else {
				messages = append(messages, "The table "+prefix+"stat_"+m+" was created (again)")
			}
		} else {
			messages = append(messages, "The table "+prefix+"stat_"+m+" was not created because it contains data.")
		}
	}

	// Inserting into base-modules
	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(ctx, MODULEID, "Statistics", VERSION)
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
