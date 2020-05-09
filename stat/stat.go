package stat

import (
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

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
