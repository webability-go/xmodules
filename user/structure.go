package user

import (
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

type StructureUser struct {
	Key  int
	Data *xdominion.XRecord
}

func CreateStructureUserByKey(sitecontext *context.Context, key int) context.Structure {
	data, _ := sitecontext.Tables["user_user"].SelectOne(key)
	if data == nil {
		return nil
	}
	return CreateStructureUserByData(sitecontext, data)
}

func CreateStructureUserByData(sitecontext *context.Context, data xdominion.XRecordDef) context.Structure {

	key, _ := data.GetInt("key")

	// Load all the data of security

	return &StructureUser{Key: key, Data: data.(*xdominion.XRecord)}
}

// ComplementData adds all the needed data from other objects /duplicable in the thread since the object will be destroyed at the end
func (sm *StructureUser) ComplementData(sitecontext *context.Context) {

}

// IsAuthorized returns true if the structure can be used on this site/language/device
func (sm *StructureUser) IsAuthorized(sitecontext *context.Context, site string, language string, device string) bool {
	return true
}

// Returns the raw data
func (sm *StructureUser) GetData() *xdominion.XRecord {
	return sm.Data
}

// Clone the whole structure
func (sm *StructureUser) Clone() context.Structure {
	cloned := &StructureUser{
		Key:  sm.Key,
		Data: sm.Data.Clone().(*xdominion.XRecord),
	}
	return cloned
}
