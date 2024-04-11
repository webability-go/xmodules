package user

import (
	"errors"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
)

func GetProfilesCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	user_profile := ds.GetTable("user_profile")
	if user_profile == nil {
		ds.Log("xmodules::user::GetCountProfilees: Error, the user_profile table is not available on this datasource")
		return 0
	}
	cnt, _ := user_profile.Count(cond)
	return cnt
}

func GetProfilesList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	user_profile := ds.GetTable("user_profile")
	if user_profile == nil {
		ds.Log("xmodules::user::GetProfileesList: Error, the user_profile table is not available on this datasource")
		return nil
	}
	data, _ := user_profile.SelectAll(cond, order, quantity, first)
	return data
}

func DeleteProfileChildren(ds applications.Datasource, skey string) error {

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

func PruneProfileChildren(ds applications.Datasource, skey string, profile string) error {

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

// GetProfileAccesses to get the list of accesses of a profile
func GetProfileAccesses(ds applications.Datasource, profile string, quantity int) (*xdominion.XRecords, error) {

	user_profileaccess := ds.GetTable("user_profileaccess")
	if user_profileaccess == nil {
		ds.Log("xmodules::user::GetProfilesOfAccess: Error, the user_profileaccess table is not available on this datasource")
		return nil, errors.New("xmodules::user::GetProfilesOfAccess: Error, the user_profileaccess table is not available on this datasource")
	}
	cond := xdominion.NewXCondition("profile", "=", profile)
	orderby := xdominion.NewXOrderBy("access", xdominion.ASC)
	data, err := user_profileaccess.SelectAll(cond, orderby, quantity)
	return data, err
}

// GetCountry to get the data of a country from cache/db in the specified language
func SetProfileAccess(ds applications.Datasource, profile string, access string, status bool) error {

	user_profileaccess := ds.GetTable("user_profileaccess")
	if user_profileaccess == nil {
		return errors.New("xmodules::user::GetProfilesOfAccess: Error, the user_profileaccess table is not available on this datasource")
	}
	cond := &xdominion.XConditions{xdominion.NewXCondition("profile", "=", profile), xdominion.NewXCondition("access", "=", access, "and")}
	data, err := user_profileaccess.SelectOne(cond)
	if err != nil {
		return err
	}
	if data != nil && !status {
		// delete entry
		_, err = user_profileaccess.Delete(cond)
	}
	if data == nil && status {
		// insert entry
		_, err = user_profileaccess.Insert(xdominion.XRecord{
			"profile": profile,
			"access":  access,
		})
	}
	return err
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetProfileUsers(ds applications.Datasource, profilekey string, quantity int) (*xdominion.XRecords, error) {

	user_userprofile := ds.GetTable("user_userprofile")
	if user_userprofile == nil {
		ds.Log("xmodules::user::GetProfilesOfAccess: Error, the user_userprofile table is not available on this datasource")
		return nil, errors.New("xmodules::user::GetProfilesOfAccess: Error, the user_userprofile table is not available on this datasource")
	}
	cond := xdominion.NewXCondition("profile", "=", profilekey)
	orderby := xdominion.NewXOrderBy("user", xdominion.ASC)
	data, err := user_userprofile.SelectAll(cond, orderby, quantity)
	return data, err
}

// HasProfilesAccessUser to get the right to an access by profiles of a user
func HasProfilesAccessUser(ds applications.Datasource, userkey int, accesskey string) (bool, error) {

	user_userprofile := ds.GetTable("user_userprofile")
	if user_userprofile == nil {
		ds.Log("xmodules::user::HasProfilesAccessUser: Error, the user_userprofile table is not available on this datasource")
		return false, errors.New("xmodules::user::HasProfilesAccessUser: Error, the user_userprofile table is not available on this datasource")
	}
	cond := xdominion.NewXCondition("user", "=", userkey)
	profiles, err := user_userprofile.SelectAll(cond)
	if err != nil {
		ds.Log("xmodules::user::HasProfilesAccessUser: " + err.Error())
		return false, err
	}
	for _, profile := range *profiles {
		profilekey, _ := profile.GetString("profile")
		hasaccess, err := HasProfileAccess(ds, profilekey, accesskey)
		if err != nil {
			return false, err
		}
		if hasaccess {
			return true, nil
		}
	}
	// no profile access found
	return false, nil
}

// HasProfileAccess to get the right to an access of a profile
func HasProfileAccess(ds applications.Datasource, profilekey string, accesskey string) (bool, error) {

	user_profileaccess := ds.GetTable("user_profileaccess")
	if user_profileaccess == nil {
		ds.Log("xmodules::user::HasProfileAccess: Error, the user_userprofile table is not available on this datasource")
		return false, errors.New("xmodules::user::GetProfilesOfAccess: Error, the user_profileaccess table is not available on this datasource")
	}
	cond := &xdominion.XConditions{xdominion.NewXCondition("profile", "=", profilekey), xdominion.NewXCondition("access", "=", accesskey, "and")}
	access, err := user_profileaccess.SelectOne(cond)
	if err != nil {
		ds.Log("xmodules::user::HasProfileAccess: " + err.Error())
		return false, err
	}
	if access != nil {
		return true, nil
	}
	return false, nil
}
