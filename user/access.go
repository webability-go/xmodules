package user

import (
	"errors"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user/assets"
)

func AddAccessGroup(ds applications.Datasource, accessgroup *assets.AccessGroup) error {

	user_accessgroup := ds.GetTable("user_accessgroup")
	if user_accessgroup == nil {
		errmsg := tools.Message(messages, "notable", "user_accessgroup", ds.GetName())
		ds.Log(errmsg)
		return errors.New(errmsg)
	}

	_, err := user_accessgroup.Upsert(accessgroup.Key, xdominion.XRecord{
		"key":         accessgroup.Key,
		"name":        accessgroup.Name,
		"description": accessgroup.Description,
	})
	if err != nil {
		ds.Log("main", tools.Message(messages, "errorupsert", "user_accessgroup", err))
		return err
	}
	return nil
}

func AddAccess(ds applications.Datasource, access *assets.Access) error {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := tools.Message(messages, "notable", "user_access", ds.GetName())
		ds.Log(errmsg)
		return errors.New(errmsg)
	}

	_, err := user_access.Upsert(access.Key, xdominion.XRecord{
		"key":         access.Key,
		"name":        access.Name,
		"group":       access.Group,
		"description": access.Description,
	})
	if err != nil {
		ds.Log("main", tools.Message(messages, "errorupsert", "user_access", err))
		return err
	}
	return nil
}

func GetCountAccesses(ds applications.Datasource, cond *xdominion.XConditions) int {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		ds.Log("xmodules::user::GetCountAccesses: Error, the user_access table is not available on this datasource")
		return 0
	}
	cnt, _ := user_access.Count()
	return cnt
}

func GetAccessesList(ds applications.Datasource, cond *xdominion.XConditions, orderby *xdominion.XOrderBy, quantity int, first int) *xdominion.XRecords {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		ds.Log("xmodules::user::GetAccessesList: Error, the user_access table is not available on this datasource")
		return nil
	}
	data, _ := user_access.SelectAll(cond, orderby, quantity, first)
	return data
}
