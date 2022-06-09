package user

import (
	"errors"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user/assets"
)

func GetAccessByQuery(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder) *xdominion.XRecord {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "GetAccessByQuery", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return nil
	}
	data, err := user_access.SelectOne(*cond, *order)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.select", "access", "GetAccessByQuery", "user_access", err))
		return nil
	}
	return data
}

func GetAccessByKey(ds applications.Datasource, key string) *xdominion.XRecord {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "GetAccessByKey", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return nil
	}
	data, err := user_access.SelectOne(key)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.select", "access", "GetAccessByKey", "user_access", err))
		return nil
	}
	return data
}

func AddAccess(ds applications.Datasource, access *assets.Access) error {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "AddAccess", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return errors.New(errmsg)
	}

	_, err := user_access.Upsert(access.Key, xdominion.XRecord{
		"key":         access.Key,
		"name":        access.Name,
		"group":       access.Group,
		"description": access.Description,
	})
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.upsert", "access", "AddAccess", "user_access", err))
	}
	return err
}

func GetAccessesCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "GetAccessesCount", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return 0
	}
	cnt, err := user_access.Count(cond)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.select", "access", "GetAccessesCount", "user_access", err))
		return 0
	}
	return cnt
}

func GetAccessesList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "GetAccessesList", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return nil
	}
	data, err := user_access.SelectAll(cond, order, quantity, first)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.select", "access", "GetAccessesList", "user_access", err))
		return nil
	}
	return data
}

func GetAccessProfiles(ds applications.Datasource, access string, quantity int) (*xdominion.XRecords, error) {

	user_profileaccess := ds.GetTable("user_profileaccess")
	if user_profileaccess == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "GetAccessProfiles", "user_profileaccess", ds.GetName())
		ds.Log("error", errmsg)
		return nil, errors.New(errmsg)
	}
	cond := xdominion.NewXCondition("access", "=", access)
	orderby := xdominion.NewXOrderBy("profile", xdominion.ASC)
	data, err := user_profileaccess.SelectAll(cond, orderby, quantity)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.select", "access", "GetAccessProfiles", "user_profileaccess", err))
		return nil, err
	}
	return data, nil
}

// GetCountry to get the data of a country from cache/db in the specified language
func GetAccessUsers(ds applications.Datasource, access string, quantity int) (*xdominion.XRecords, error) {

	user_useraccess := ds.GetTable("user_useraccess")
	if user_useraccess == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "GetAccessUsers", "user_useraccess", ds.GetName())
		ds.Log("error", errmsg)
		return nil, errors.New(errmsg)
	}
	cond := xdominion.NewXCondition("access", "=", access)
	orderby := xdominion.NewXOrderBy("user", xdominion.ASC)
	data, err := user_useraccess.SelectAll(cond, orderby, quantity)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.select", "access", "GetAccessUsers", "user_useraccess", err))
		return nil, err
	}
	return data, nil
}

func DeleteAccessChildren(ds applications.Datasource, skey string) error {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "DeleteAccessChildren", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return errors.New(errmsg)
	}
	//	_, err := user_access.Delete(xdominion.NewXCondition("group", "=", skey))
	//	if err != nil {
	//		ds.Log("error", tools.Message(messages, "error.delete", "accessgroup", "DeleteAccessGroupChildren", "user_access", err))
	//	}
	return nil
}

func PruneAccessChildren(ds applications.Datasource, skey string, group string) error {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "PruneAccessChildren", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return errors.New(errmsg)
	}
	//	_, err := user_access.Update(xdominion.NewXCondition("group", "=", skey), xdominion.XRecord{"group": group})
	//	if err != nil {
	//		ds.Log("error", tools.Message(messages, "error.delete", "accessgroup", "PruneAccessGroupChildren", "user_access", err))
	//	}
	return nil
}
