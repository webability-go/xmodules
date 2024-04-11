package translation

import (
	//  "fmt"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
)

func GetSourceByKey(ds applications.Datasource, key int) *xdominion.XRecord {

	kl_traduccionfuente := ds.GetTable("kl_traduccionfuente")
	if kl_traduccionfuente == nil {
		ds.Log("xmodules::translation::GetSourceByKey: Error, the kl_traduccionfuente table is not available on this datasource")
		return nil
	}
	data, _ := kl_traduccionfuente.SelectOne(key)
	return data
}

func GetSourcesCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	kl_traduccionfuente := ds.GetTable("kl_traduccionfuente")
	if kl_traduccionfuente == nil {
		ds.Log("xmodules::translation::GetSourcesCount: Error, the kl_traduccionfuente table is not available on this datasource")
		return 0
	}
	cnt, _ := kl_traduccionfuente.Count(cond)
	return cnt
}

func GetSourcesList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	kl_traduccionfuente := ds.GetTable("kl_traduccionfuente")
	if kl_traduccionfuente == nil {
		ds.Log("xmodules::translation::GetSourcesList: Error, the kl_traduccionfuente table is not available on this datasource")
		return nil
	}
	data, _ := kl_traduccionfuente.SelectAll(cond, order, quantity, first)
	return data
}

func DeleteSourceChildren(ds applications.Datasource, skey string) error {

	return nil
}

func PruneSourceChildren(ds applications.Datasource, skey string, channel string) error {

	return nil
}
