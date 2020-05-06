package user

import (
	"fmt"
	"strconv"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/context"
)

// GetCountry to get the data of a country from cache/db in the specified language
func GetUser(sitecontext *context.Context, key int) *StructureUser {

	cache := sitecontext.GetCache("user:users")
	if cache == nil {
		sitecontext.Log("main", "xmodules::user::GetUser: Error, the user cache is not available on this site context")
		return nil
	}

	data, _ := cache.Get(strconv.Itoa(key))
	if data == nil {
		sm := CreateStructureUserByKey(sitecontext, key)
		if sm == nil {
			sitecontext.Log("graph", "xmodules::user::GetUser: There is no user created: ", key)
			return nil
		}
		cache.Set(strconv.Itoa(key), sm)
		return sm.(*StructureUser)
	}
	return data.(*StructureUser)
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetUsersList(sitecontext *context.Context) *xdominion.XRecords {

	user_user := sitecontext.GetTable("user_user")
	if user_user == nil {
		sitecontext.Log("xmodules::user::GetUsersList: Error, the user_user table is not available on this context")
		return nil
	}
	data, _ := user_user.SelectAll()
	return data
}

func GetUserByCredentials(sitecontext *context.Context, username string, password string) *StructureUser {

	user_user := sitecontext.GetTable("user_user")
	if user_user == nil {
		sitecontext.Log("xmodules::user::GetUserByCredentials: Error, the user_user table is not available on this context")
		return nil
	}
	data, _ := user_user.SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("un", "=", username),
		xdominion.NewXCondition("pw", "=", password, "and"),
		xdominion.NewXCondition("status", "!=", "X", "and"),
	})
	if data == nil {
		return nil
	}
	sm := CreateStructureUserByData(sitecontext, data)
	if sm == nil {
		sitecontext.Log("graph", "xmodules::user::GetUser: There is no user created: ", fmt.Sprint(data))
		return nil
	}
	return sm.(*StructureUser)
}
