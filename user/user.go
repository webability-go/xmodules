package user

import (
	"fmt"
	"strconv"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/base"
)

// GetCountry to get the data of a country from cache/db in the specified language
func GetUser(ds *base.Datasource, key int) *StructureUser {

	cache := ds.GetCache("user:users")
	if cache == nil {
		ds.Log("main", "xmodules::user::GetUser: Error, the user cache is not available on this site datasource")
		return nil
	}

	data, _ := cache.Get(strconv.Itoa(key))
	if data == nil {
		sm := CreateStructureUserByKey(ds, key)
		if sm == nil {
			ds.Log("graph", "xmodules::user::GetUser: There is no user created: ", key)
			return nil
		}
		cache.Set(strconv.Itoa(key), sm)
		return sm.(*StructureUser)
	}
	return data.(*StructureUser)
}

// GetUsersList to get a list of all the users
func GetUsersList(ds *base.Datasource) *xdominion.XRecords {

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		ds.Log("xmodules::user::GetUsersList: Error, the user_user table is not available on this datasource")
		return nil
	}
	data, _ := user_user.SelectAll()
	return data
}

func GetUserByCredentials(ds *base.Datasource, username string, password string) *StructureUser {

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		ds.Log("xmodules::user::GetUserByCredentials: Error, the user_user table is not available on this datasource")
		return nil
	}
	data, err := user_user.SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("un", "=", username),
		xdominion.NewXCondition("pw", "=", password, "and"),
		xdominion.NewXCondition("status", "!=", "X", "and"),
	})
	if err != nil {
		ds.Log("xmodules::user::GetUserByCredentials:" + err.Error())
	}
	if data == nil {
		return nil
	}
	sm := CreateStructureUserByData(ds, data)
	if sm == nil {
		ds.Log("graph", "xmodules::user::GetUser: There is no user created: ", fmt.Sprint(data))
		return nil
	}
	return sm.(*StructureUser)
}
