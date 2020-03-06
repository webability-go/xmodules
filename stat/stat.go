package stat

import (
	"strconv"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "stat"
	VERSION  = "2.0.0"
)

func InitModule(sitecontext *context.Context, prefix string, databasename string) error {

	buildTables(sitecontext, prefix, databasename)
	sitecontext.SetModule(MODULEID, VERSION)

	return nil
}

func SynchronizeModule(sitecontext *context.Context, prefix string) []string {

	messages := []string{}

	// Needed modules: context and translation
	vc := context.ModuleInstalledVersion(sitecontext, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/country.")
		return messages
	}

	for i := 1; i < 13; i++ {
		m := ""
		if i < 10 {
			m = "0"
		}
		m += strconv.Itoa(i)

		messages = append(messages, "Analysing "+prefix+"stat_"+m+" table.")
		table := sitecontext.GetTable(prefix + "stat_" + m)
		if table == nil {
			messages = append(messages, "xmodules::stat::SynchronizeModule: Error, the table does not exist in the context: "+prefix+"stat_"+m)
			return messages
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
	err := context.AddModule(sitecontext, MODULEID, "Statistics", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages
}

func RegisterStat(sitecontext *context.Context, prefix string, data xdominion.XRecord) {

	table := sitecontext.GetTable(prefix + "stat_" + getMonth())
	if table == nil {
		sitecontext.Log("main", "xmodules::stat::RegisterStat: Error, the table does not exist in the context: ", prefix+"stat_"+getMonth())
		return
	}

	data.Set("clave", 0)
	_, err := table.Insert(data)
	if err != nil {
		sitecontext.Log("main", "Error insertando el log:", err)
	}
}

// TODO(phil) hacer las funciones RegisterMiss y RegisterSys
