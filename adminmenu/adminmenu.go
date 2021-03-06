package adminmenu

import (
	"errors"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
)

func AddGroup(ds applications.Datasource, key string, name string) error {
	adminmenu_group := ds.GetTable("adminmenu_group")
	if adminmenu_group == nil {
		errmessage := "xmodules::adminmenu::AddGroup: Error, the adminmenu_group table is not available on this datasource"
		ds.Log(errmessage)
		return errors.New(errmessage)
	}
	_, err := adminmenu_group.Upsert(key, xdominion.XRecord{
		"key":  key,
		"name": name,
	})
	if err != nil {
		ds.Log("main", "Error inserting in adminmenu_group", err)
		return err
	}
	return nil
}

func GetGroup(ds applications.Datasource, key string) (*xdominion.XRecord, error) {
	return nil, nil
}

func AddOption(ds applications.Datasource, data *xdominion.XRecord) error {

	adminmenu_option := ds.GetTable("adminmenu_option")
	if adminmenu_option == nil {
		errmessage := "xmodules::adminmenu::AddOption: Error, the adminmenu_option table is not available on this datasource"
		ds.Log(errmessage)
		return errors.New(errmessage)
	}

	// Verify if it exists already

	key, _ := data.GetString("key")
	_, err := adminmenu_option.Upsert(key, data)
	if err != nil {
		ds.Log("main", "Error upserting in adminmenu_option", err)
		return err
	}
	return nil
}

func GetOption(ds applications.Datasource, key string) (*xdominion.XRecord, error) {
	return nil, nil
}

func GetMenu(ds applications.Datasource, group string, father string) (*xdominion.XRecords, error) {

	adminmenu_option := ds.GetTable("adminmenu_option")
	if adminmenu_option == nil {
		errmessage := "xmodules::adminmenu::AddOption: Error, the adminmenu_option table is not available on this datasource"
		ds.Log(errmessage)
		return nil, errors.New(errmessage)
	}

	var sfather interface{}
	if father != "" {
		sfather = father
	}
	cond := xdominion.XConditions{
		xdominion.NewXCondition("group", "=", group),
		xdominion.NewXCondition("father", "=", sfather, "and"),
	}
	data, err := adminmenu_option.SelectAll(cond)
	if err != nil {
		ds.Log("xmodules::user::GetUserByCredentials:" + err.Error())
		return nil, err
	}

	return data, nil
}
