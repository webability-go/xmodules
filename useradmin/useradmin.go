package useradmin

import (
	"github.com/webability-go/xmodules/base"

	"github.com/webability-go/xdominion"
)

// GetUsersList to get a list of all the users
func GetOptions(ds *base.Datasource, group string, father int) *xdominion.XRecords {

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		ds.Log("xmodules::user::GetUsersList: Error, the user_user table is not available on this datasource")
		return nil
	}
	data, _ := user_user.SelectAll()
	return data
}
