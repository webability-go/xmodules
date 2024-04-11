package user

import (
	"errors"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user/assets"
)

// AddAccessGroup is generally used by xmodules installers
func AddAccessGroup(ds applications.Datasource, accessgroup *assets.AccessGroup) error {

	user_accessgroup := ds.GetTable("user_accessgroup")
	if user_accessgroup == nil {
		errmsg := tools.Message(messages, "error.notable", "accessgroup", "AddAccessGroup", "user_accessgroup", ds.GetName())
		ds.Log("error", errmsg)
		return errors.New(errmsg)
	}

	_, err := user_accessgroup.Upsert(accessgroup.Key, xdominion.XRecord{
		"key":         accessgroup.Key,
		"name":        accessgroup.Name,
		"description": accessgroup.Description,
	})
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.upsert", "accessgroup", "AddAccessGroup", "user_accessgroup", err))
	}
	return err
}

func DelAccessGroupByKey(ds applications.Datasource, key string) error {

	user_accessgroup := ds.GetTable("user_accessgroup")
	if user_accessgroup == nil {
		errmsg := tools.Message(messages, "error.notable", "access", "DelAccessGroup", "user_accessgroup", ds.GetName())
		ds.Log("error", errmsg)
		return errors.New(errmsg)
	}

	_, err := user_accessgroup.Delete(key)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.delete", "access", "DelAccessGroup", "user_accessgroup", err))
	}
	return err
}

func GetAccessGroupsCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	user_accessgroup := ds.GetTable("user_accessgroup")
	if user_accessgroup == nil {
		errmsg := tools.Message(messages, "error.notable", "accessgroup", "GetAccessGroupsCount", "user_accessgroup", ds.GetName())
		ds.Log("error", errmsg)
		return 0
	}
	cnt, err := user_accessgroup.Count(cond)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.select", "accessgroup", "GetAccessGroupsCount", "user_accessgroup", err))
		return 0
	}
	return cnt
}

func GetAccessGroupsList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	user_accessgroup := ds.GetTable("user_accessgroup")
	if user_accessgroup == nil {
		errmsg := tools.Message(messages, "error.notable", "accessgroup", "GetAccessGroupsList", "user_accessgroup", ds.GetName())
		ds.Log("error", errmsg)
		return nil
	}
	data, err := user_accessgroup.SelectAll(cond, order, quantity, first)
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.select", "accessgroup", "GetAccessGroupsList", "user_accessgroup", err))
		return nil
	}
	return data
}

func DeleteAccessGroupChildren(ds applications.Datasource, skey string) error {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "accessgroup", "DeleteAccessGroupChildren", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return errors.New(errmsg)
	}
	accesses, err := user_access.SelectAll(xdominion.NewXCondition("group", "=", skey))
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.delete", "accessgroup", "DeleteAccessGroupChildren", "user_access", err))
		return err
	}
	if accesses == nil {
		return nil
	}
	for _, r := range *accesses {
		akey, _ := r.GetString("key")
		err = DeleteAccessChildren(ds, akey)
		if err != nil {
			ds.Log("error", tools.Message(messages, "error.delete", "accessgroup", "DeleteAccessGroupChildren", "user_access", err))
			return err
		}
	}
	_, err = user_access.Delete(xdominion.NewXCondition("group", "=", skey))
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.delete", "accessgroup", "DeleteAccessGroupChildren", "user_access", err))
	}
	return err
}

func PruneAccessGroupChildren(ds applications.Datasource, skey string, group string) error {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "error.notable", "accessgroup", "PruneAccessGroupChildren", "user_access", ds.GetName())
		ds.Log("error", errmsg)
		return errors.New(errmsg)
	}
	_, err := user_access.Update(xdominion.NewXCondition("group", "=", skey), xdominion.XRecord{"group": group})
	if err != nil {
		ds.Log("error", tools.Message(messages, "error.delete", "accessgroup", "PruneAccessGroupChildren", "user_access", err))
	}
	return err
}
