package user

import (
	"strings"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
)

func SetUserParam(ds applications.Datasource, user int, param string, value interface{}) {

	user_parameter := ds.GetTable("user_parameter")
	if user_parameter == nil {
		ds.Log("xmodules::user::SetUserParam: Error, the user_parameter table is not available on this datasource")
		return
	}
	data, _ := user_parameter.SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("id", "=", param),
		xdominion.NewXCondition("user", "=", user, "and"),
	})
	if data == nil {
		_, err := user_parameter.Insert(xdominion.XRecord{
			"key":   0,
			"user":  user,
			"id":    param,
			"value": value,
		})
		if err != nil {
			ds.Log("xmodules::user::SetUserParam: Error inserting in the user_parameter table", err)
		}
		return
	}
	key, _ := data.GetInt("key")
	_, err := user_parameter.Update(key, xdominion.XRecord{
		"value": value,
	})
	if err != nil {
		ds.Log("xmodules::user::SetUserParam: Error", err)
	}
}

func AddUserParam(ds applications.Datasource, user int, param string, value interface{}) {

	user_parameter := ds.GetTable("user_parameter")
	if user_parameter == nil {
		ds.Log("xmodules::user::SetUserParam: Error, the user_parameter table is not available on this datasource")
		return
	}
	_, err := user_parameter.Insert(xdominion.XRecord{
		"key":   0,
		"user":  user,
		"id":    param,
		"value": value,
	})
	if err != nil {
		ds.Log("xmodules::user::AddUserParam: Error inserting in the user_parameter table", err)
	}
}

func GetUserParam(ds applications.Datasource, user int, param string) string {

	user_parameter := ds.GetTable("user_parameter")
	if user_parameter == nil {
		ds.Log("xmodules::user::SetUserParam: Error, the user_parameter table is not available on this datasource")
		return ""
	}
	data, _ := user_parameter.SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("id", "=", param),
		xdominion.NewXCondition("user", "=", user, "and"),
	})
	if data == nil {
		return ""
	}

	value, _ := data.GetString("value")
	return strings.TrimSpace(value)
}

func DelUserParam(ds applications.Datasource, user int, param string) {

	user_parameter := ds.GetTable("user_parameter")
	if user_parameter == nil {
		ds.Log("xmodules::user::SetUserParam: Error, the user_parameter table is not available on this datasource")
		return
	}
	_, err := user_parameter.Delete(xdominion.XConditions{
		xdominion.NewXCondition("id", "=", param),
		xdominion.NewXCondition("user", "=", user, "and"),
	})
	if err != nil {
		ds.Log("xmodules::user::DelUserParam: Error", err)
	}
}
