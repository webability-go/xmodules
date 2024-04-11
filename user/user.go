package user

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xdominion"
)

// GetCountry to get the data of a country from cache/db in the specified language
func GetUserByKey(ds applications.Datasource, key int) *xdominion.XRecord {

	cache := ds.GetCache("user:users")
	if cache == nil {
		ds.Log("main", "xmodules::user::GetUserByKey: Error, the user cache is not available on this site datasource")
		return nil
	}

	data, _ := cache.Get(strconv.Itoa(key))
	if data == nil {
		sm := CreateStructureUserByKey(ds, key)
		if sm == nil {
			ds.Log("graph", "xmodules::user::GetUserByKey: There is no user created: ", key)
			return nil
		}
		cache.Set(strconv.Itoa(key), sm)
		return sm.(*StructureUser).Data
	}
	return data.(*StructureUser).Data
}

func GetUserByCredentials(ds applications.Datasource, username string, password string) *StructureUser {

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

func GetUsersCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		ds.Log("xmodules::user::GetCountProfilees: Error, the user_user table is not available on this datasource")
		return 0
	}
	cnt, _ := user_user.Count(cond)
	return cnt
}

func GetUsersList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		ds.Log("xmodules::user::GetProfileesList: Error, the user_user table is not available on this datasource")
		return nil
	}
	data, _ := user_user.SelectAll(cond, order, quantity, first)
	return data
}

func DeleteUserChildren(ds applications.Datasource, key int) error {

	/*
		user_profile := ds.GetTable("user_profile")
		if user_profile == nil {
			errmsg := tools.Message(messages, "notable", "user_profile", ds.GetName())
			ds.Log(errmsg)
			return errors.New(errmsg)
		}
		_, err := user_profile.Delete(xdominion.NewXCondition("group", "=", skey))
		return err
	*/
	return nil
}

func PruneUserChildren(ds applications.Datasource, key int, user int) error {

	/*
		user_profile := ds.GetTable("user_profile")
		if user_profile == nil {
			errmsg := tools.Message(messages, "notable", "user_profile", ds.GetName())
			ds.Log(errmsg)
			return errors.New(errmsg)
		}
		_, err := user_profile.Update(xdominion.NewXCondition("group", "=", skey), xdominion.XRecord{"group": group})
		return err
	*/
	return nil
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetUserAccessByKeys(ds applications.Datasource, userkey int, accesskey string) *xdominion.XRecord {

	return nil
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetUserAccesses(ds applications.Datasource, userkey int, quantity int) (*xdominion.XRecords, error) {

	user_useraccess := ds.GetTable("user_useraccess")
	if user_useraccess == nil {
		ds.Log("xmodules::user::GetProfilesOfAccess: Error, the user_useraccess table is not available on this datasource")
		return nil, errors.New("xmodules::user::GetProfilesOfAccess: Error, the user_userprofile table is not available on this datasource")
	}
	cond := xdominion.NewXCondition("user", "=", userkey)
	orderby := xdominion.NewXOrderBy("access", xdominion.ASC)
	data, err := user_useraccess.SelectAll(cond, orderby, quantity)
	return data, err
}

// GetUserAccess to get the record of an access for the user
func GetUserAccess(ds applications.Datasource, userkey int, access string) (*xdominion.XRecord, error) {

	user_useraccess := ds.GetTable("user_useraccess")
	if user_useraccess == nil {
		ds.Log("xmodules::user::GetUserAccess: Error, the user_useraccess table is not available on this datasource")
		return nil, errors.New("xmodules::user::GetUserAccess: Error, the user_userprofile table is not available on this datasource")
	}
	cond := &xdominion.XConditions{xdominion.NewXCondition("user", "=", userkey), xdominion.NewXCondition("access", "=", access, "and")}
	data, err := user_useraccess.SelectOne(cond)
	return data, err
}

// GetCountry to get the data of a country from cache/db in the specified language
func SetUserAccess(ds applications.Datasource, user int, access string, status int) error {

	user_useraccess := ds.GetTable("user_useraccess")
	if user_useraccess == nil {
		return errors.New("xmodules::user::GetProfilesOfAccess: Error, the user_useraccess table is not available on this datasource")
	}
	cond := xdominion.XConditions{xdominion.NewXCondition("user", "=", user), xdominion.NewXCondition("access", "=", access, "and")}
	data, err := user_useraccess.SelectOne(cond)
	if err != nil {
		return err
	}
	if data != nil && status == -1 {
		// delete entry
		_, err = user_useraccess.Delete(cond)
	}
	if data == nil && status != -1 {
		// insert entry
		_, err = user_useraccess.Insert(xdominion.XRecord{
			"user":   user,
			"access": access,
			"denied": status,
		})
	}
	if data != nil && status != -1 {
		// insert entry
		_, err = user_useraccess.Update(cond, xdominion.XRecord{
			"denied": status,
		})
	}
	return err
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetUserProfiles(ds applications.Datasource, userkey int, quantity int) (*xdominion.XRecords, error) {

	user_userprofile := ds.GetTable("user_userprofile")
	if user_userprofile == nil {
		ds.Log("xmodules::user::GetProfilesOfAccess: Error, the user_userprofile table is not available on this datasource")
		return nil, errors.New("xmodules::user::GetProfilesOfAccess: Error, the user_userprofile table is not available on this datasource")
	}
	cond := xdominion.NewXCondition("user", "=", userkey)
	orderby := xdominion.NewXOrderBy("profile", xdominion.ASC)
	data, err := user_userprofile.SelectAll(cond, orderby, quantity)
	return data, err
}

// GetCountry to get the data of a country from cache/db in the specified language
func SetUserProfile(ds applications.Datasource, user int, profile string, status bool) error {

	user_userprofile := ds.GetTable("user_userprofile")
	if user_userprofile == nil {
		return errors.New("xmodules::user::GetProfilesOfAccess: Error, the user_userprofile table is not available on this datasource")
	}
	cond := xdominion.XConditions{xdominion.NewXCondition("user", "=", user), xdominion.NewXCondition("profile", "=", profile, "and")}
	data, err := user_userprofile.SelectOne(cond)
	if err != nil {
		return err
	}
	if data != nil && !status {
		// delete entry
		_, err = user_userprofile.Delete(cond)
	}
	if data == nil && status {
		// insert entry
		_, err = user_userprofile.Insert(xdominion.XRecord{
			"user":    user,
			"profile": profile,
		})
	}
	return err
}
