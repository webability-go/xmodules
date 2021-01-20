package user

import (
	"errors"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/user/assets"
)

func AddAccess(ds *base.Datasource, access *assets.Access) error {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		errmsg := "xmodules::user::AddAccess: Error, the user_access table is not available on this datasource " + ds.Name
		ds.Log(errmsg)
		return errors.New(errmsg)
	}
	return nil
}

func GetCountAccesses(ds *base.Datasource, cond *xdominion.XConditions) int {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		ds.Log("xmodules::user::GetCountAccesses: Error, the user_access table is not available on this datasource")
		return 0
	}
	cnt, _ := user_access.Count()
	return cnt
}

func GetAccessesList(ds *base.Datasource, cond *xdominion.XConditions, orderby *xdominion.XOrderBy, quantity int, first int) *xdominion.XRecords {

	user_access := ds.GetTable("user_access")
	if user_access == nil {
		ds.Log("xmodules::user::GetAccessesList: Error, the user_access table is not available on this datasource")
		return nil
	}
	data, _ := user_access.SelectAll(cond, orderby, quantity, first)
	return data
}
