package assets

import (
	"log"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID   = "user"
	VERSION    = "0.0.1"
	DATASOURCE = "userdatasource"
)

var Needs = []string{"base"}

type ModuleEntries struct {
	// Access Groups
	GetAccessGroupsCount      func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetAccessGroupsList       func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	DeleteAccessGroupChildren func(ds applications.Datasource, skey string) error
	PruneAccessGroupChildren  func(ds applications.Datasource, skey string, group string) error

	// Accesses
	GetAccessByKey func(ds applications.Datasource, key string) *xdominion.XRecord
	//	GetAccessByQuery     func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder) *xdominion.XRecord
	GetAccessesCount     func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetAccessesList      func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	DeleteAccessChildren func(ds applications.Datasource, skey string) error
	PruneAccessChildren  func(ds applications.Datasource, skey string, access string) error
	GetAccessProfiles    func(ds applications.Datasource, key string, quantity int) (*xdominion.XRecords, error)
	GetAccessUsers       func(ds applications.Datasource, key string, quantity int) (*xdominion.XRecords, error)

	// profiles
	GetProfilesCount      func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetProfilesList       func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	DeleteProfileChildren func(ds applications.Datasource, skey string) error
	PruneProfileChildren  func(ds applications.Datasource, skey string, profile string) error
	GetProfileAccesses    func(ds applications.Datasource, key string, quantity int) (*xdominion.XRecords, error)
	SetProfileAccess      func(ds applications.Datasource, key string, access string, status bool) error
	GetProfileUsers       func(ds applications.Datasource, key string, quantity int) (*xdominion.XRecords, error)

	// users
	GetUserByKey       func(ds applications.Datasource, key int) *xdominion.XRecord
	GetUsersCount      func(ds applications.Datasource, cond *xdominion.XConditions) int
	GetUsersList       func(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords
	DeleteUserChildren func(ds applications.Datasource, key int) error
	PruneUserChildren  func(ds applications.Datasource, key int, user int) error
	GetUserAccesses    func(ds applications.Datasource, key int, quantity int) (*xdominion.XRecords, error)
	SetUserAccess      func(ds applications.Datasource, key int, access string, status int) error
	GetUserProfiles    func(ds applications.Datasource, key int, quantity int) (*xdominion.XRecords, error)
	SetUserProfile     func(ds applications.Datasource, key int, profile string, status bool) error

	// User Params
	SetUserParam func(ds applications.Datasource, user int, param string, value interface{})
	AddUserParam func(ds applications.Datasource, user int, param string, value interface{})
	GetUserParam func(ds applications.Datasource, user int, param string) string
	DelUserParam func(ds applications.Datasource, user int, param string)

	// security
	HasAccess func(ds applications.Datasource, userid int, args ...interface{}) bool
}

func GetEntries(logger *log.Logger) *ModuleEntries {
	me := base.GetEntries(logger, MODULEID)
	if me == nil {
		return nil
	}
	lme, ok := me.(*ModuleEntries)
	if !ok {
		return nil
	}
	return lme
}
