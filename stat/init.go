// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package stat

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "stat"
	VERSION  = "0.0.1"
)

func init() {
	m := &context.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Stat logs", language.Spanish: "Bitacoras", language.French: "Registres"},
		Needs:        []string{"context"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	context.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to ctx::database
func Setup(ctx *context.Context, prefix string) ([]string, error) {

	buildTables(ctx, prefix)
	//	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	//	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ctx *context.Context, prefix string) ([]string, error) {

	messages := []string{}

	// Needed modules: context and translation
	vc := context.ModuleInstalledVersion(ctx, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/stat.")
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
			messages = append(messages, "xmodules::stat::SynchronizeModule: Error, the table does not exist in the context: "+prefix+"stat_"+m)
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

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err := context.AddModule(ctx, MODULEID, "Statistics", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages, nil
}
