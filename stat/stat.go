package stat

import (
	"strconv"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "stat"
	VERSION  = "1.0.0"
)

func InitModule(sitecontext *context.Context, prefix string, databasename string) error {

	buildTables(sitecontext, prefix, databasename)
	sitecontext.Modules[MODULEID] = VERSION

	return nil
}

func SynchronizeModule(sitecontext *context.Context, prefix string) []string {

	messages := []string{}
	for i := 1; i < 13; i++ {
		m := ""
		if i < 10 {
			m = "0"
		}
		m += strconv.Itoa(i)

		messages = append(messages, "Analysing "+prefix+"stat_"+m+" table.")
		num, err := sitecontext.Tables[prefix+"stat_"+m].Count(nil)
		if err != nil || num == 0 {
			err1 := sitecontext.Tables[prefix+"stat_"+m].Synchronize()
			if err1 != nil {
				messages = append(messages, "The table "+prefix+"stat_"+m+" was not created: "+err1.Error())
			} else {
				messages = append(messages, "The table "+prefix+"stat_"+m+" was created (again)")
			}
		} else {
			messages = append(messages, "The table "+prefix+"stat_"+m+" was not created because it contains data.")
		}
	}
	return messages
}

func RegisterStat(sitecontext *context.Context, prefix string, data xdominion.XRecord) {
	data.Set("clave", 0)
	_, err := sitecontext.Tables[prefix+"stat_"+getMonth()].Insert(data)
	if err != nil {
		sitecontext.Logs["main"].Println("Error insertando el log:", err)
	}
}

// TODO(phil) hacer las funciones RegisterMiss y RegisterSys
