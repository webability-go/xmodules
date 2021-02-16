package user

import (
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/base"
)

type StructureUser struct {
	Key  int
	Data *xdominion.XRecord
}

func CreateStructureUserByKey(ds applications.Datasource, key int) base.Structure {

	user_user := ds.GetTable("user_user")
	if user_user == nil {
		ds.Log("xmodules::user::CreateStructureUserByKey: Error, the user_user table is not available on this datasource")
		return nil
	}

	data, _ := user_user.SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureUserByData(ds, data)
}

func CreateStructureUserByData(ds applications.Datasource, data xdominion.XRecordDef) base.Structure {

	key, _ := data.GetInt("key")

	// Load all the data of security

	return &StructureUser{Key: key, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureUser) ComplementData(ds applications.Datasource) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureUser) IsAuthorized(ds applications.Datasource, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructureUser) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructureUser) Clone() base.Structure {
	cloned := &StructureUser{
		Key:  sm.Key,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
