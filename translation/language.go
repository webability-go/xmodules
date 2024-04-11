package translation

import (
	//  "fmt"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
)

func GetLanguageByKey(ds applications.Datasource, key string) *xdominion.XRecord {

	kl_language := ds.GetTable("kl_language")
	if kl_language == nil {
		ds.Log("xmodules::translation::GetLanguageByKey: Error, the kl_language table is not available on this datasource")
		return nil
	}
	data, _ := kl_language.SelectOne(key)
	return data
}

func GetLanguagesCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	kl_language := ds.GetTable("kl_language")
	if kl_language == nil {
		ds.Log("xmodules::video::GetCountProfilees: Error, the kl_language table is not available on this datasource")
		return 0
	}
	cnt, _ := kl_language.Count(cond)
	return cnt
}

func GetLanguagesList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	kl_language := ds.GetTable("kl_language")
	if kl_language == nil {
		ds.Log("xmodules::video::GetProfileesList: Error, the kl_language table is not available on this datasource")
		return nil
	}
	data, _ := kl_language.SelectAll(cond, order, quantity, first)
	return data
}

func DeleteLanguageChildren(ds applications.Datasource, skey string) error {

	return nil
}

func PruneLanguageChildren(ds applications.Datasource, skey string, channel string) error {

	return nil
}
