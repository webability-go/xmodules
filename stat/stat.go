package stat

import (
	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xdominion"
)

func RegisterStat(ds applications.Datasource, prefix string, data xdominion.XRecord) {

	table := ds.GetTable(prefix + "stat_" + getMonth())
	if table == nil {
		ds.Log("main", "xmodules::stat::RegisterStat: Error, the table does not exist in the context: ", prefix+"stat_"+getMonth())
		return
	}

	data.Set("clave", 0)
	_, err := table.Insert(data)
	if err != nil {
		ds.Log("main", "Error insertando el log:", err)
	}
}
