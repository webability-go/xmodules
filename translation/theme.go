package translation

import (
	//  "fmt"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/tools"

	"github.com/webability-go/xamboo/applications"
)

func GetThemeByKey(ds applications.Datasource, key int) *xdominion.XRecord {

	kl_traducciontema := ds.GetTable("kl_traducciontema")
	if kl_traducciontema == nil {
		ds.Log("xmodules::translation::GetThemeByKey: Error, the kl_traducciontema table is not available on this datasource")
		return nil
	}
	data, _ := kl_traducciontema.SelectOne(key)
	return data
}

func GetThemeByName(ds applications.Datasource, name string) *xdominion.XRecord {

	kl_traducciontema := ds.GetTable("kl_traducciontema")
	if kl_traducciontema == nil {
		ds.Log("xmodules::translation::GetThemeByName: Error, the kl_traducciontema table is not available on this datasource")
		return nil
	}
	cond := xdominion.XConditions{xdominion.NewXCondition("nombre", "=", name)}
	data, _ := kl_traducciontema.SelectOne(&cond)
	return data
}

// AddThemevale, le encargo a  is generally used by xmodules installers
func AddTheme(ds applications.Datasource, data *xdominion.XRecord) error {

	kl_traducciontema := ds.GetTable("kl_traducciontema")
	if kl_traducciontema == nil {
		ds.Log("xmodules::translation::AddTheme: Error, the kl_traducciontema table is not available on this datasource")
		return nil
	}

	key, _ := data.GetInt("clave")
	_, err := kl_traducciontema.Upsert(key, data)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.upsert", "translation", "AddTheme", "kl_traducciontema", err))
	}
	return err
}

func DelThemeByKey(ds applications.Datasource, key int) error {

	kl_traducciontema := ds.GetTable("kl_traducciontema")
	if kl_traducciontema == nil {
		ds.Log("xmodules::translation::DelThemeByKey: Error, the kl_traducciontema table is not available on this datasource")
		return nil
	}
	_, err := kl_traducciontema.Delete(key)
	return err
}

func GetThemesCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	kl_traducciontema := ds.GetTable("kl_traducciontema")
	if kl_traducciontema == nil {
		ds.Log("xmodules::video::GetCountProfilees: Error, the kl_traducciontema table is not available on this datasource")
		return 0
	}
	cnt, _ := kl_traducciontema.Count(cond)
	return cnt
}

func GetThemesList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	kl_traducciontema := ds.GetTable("kl_traducciontema")
	if kl_traducciontema == nil {
		ds.Log("xmodules::video::GetProfileesList: Error, the kl_traducciontema table is not available on this datasource")
		return nil
	}
	data, _ := kl_traducciontema.SelectAll(cond, order, quantity, first)
	return data
}

func DeleteThemeChildren(ds applications.Datasource, skey string) error {

	return nil
}

func PruneThemeChildren(ds applications.Datasource, skey string, channel string) error {

	return nil
}
