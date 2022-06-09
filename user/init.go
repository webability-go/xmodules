// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package user

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user/assets"
)

//go:embed languages/*.language
var fsmessages embed.FS
var messages *map[language.Tag]*xcore.XLanguage

func init() {
	messages = tools.BuildMessagesFS(fsmessages, "languages")
	m := &base.Module{
		ID:      assets.MODULEID,
		Version: assets.VERSION,
		Languages: map[language.Tag]string{
			language.English: tools.Message(messages, "MODULENAME", language.English),
			language.Spanish: tools.Message(messages, "MODULENAME", language.Spanish),
			language.French:  tools.Message(messages, "MODULENAME", language.French),
		},
		Needs:         assets.Needs,
		FSetup:        Setup,        // Called once at the main system startup, once PER CREATED xmodule CONTEXT (if set)
		FSynchronize:  Synchronize,  // Called only to create/rebuild database objects and others on demand (if set)
		FStartContext: StartContext, // Called each time a new Server context is created  (if set)
		Entries: &assets.ModuleEntries{
			// access groups
			GetAccessGroupsCount:      GetAccessGroupsCount,
			GetAccessGroupsList:       GetAccessGroupsList,
			DeleteAccessGroupChildren: DeleteAccessGroupChildren,
			PruneAccessGroupChildren:  PruneAccessGroupChildren,

			// accesses
			GetAccessByKey: GetAccessByKey,
			//	GetAccessByQuery:     GetAccessByQuery,
			GetAccessesCount:     GetAccessesCount,
			GetAccessesList:      GetAccessesList,
			DeleteAccessChildren: DeleteAccessChildren,
			PruneAccessChildren:  PruneAccessChildren,
			GetAccessUsers:       GetAccessUsers,
			GetAccessProfiles:    GetAccessProfiles,

			// profiles
			GetProfilesCount:      GetProfilesCount,
			GetProfilesList:       GetProfilesList,
			DeleteProfileChildren: DeleteProfileChildren,
			PruneProfileChildren:  PruneProfileChildren,
			GetProfileAccesses:    GetProfileAccesses,
			SetProfileAccess:      SetProfileAccess,
			GetProfileUsers:       GetProfileUsers,

			// users
			GetUserByKey:       GetUserByKey,
			GetUsersCount:      GetUsersCount,
			GetUsersList:       GetUsersList,
			DeleteUserChildren: DeleteUserChildren,
			PruneUserChildren:  PruneUserChildren,
			GetUserAccesses:    GetUserAccesses,
			SetUserAccess:      SetUserAccess,
			GetUserProfiles:    GetUserProfiles,
			SetUserProfile:     SetUserProfile,

			// Params
			SetUserParam: SetUserParam,
			AddUserParam: AddUserParam,
			GetUserParam: GetUserParam,
			DelUserParam: DelUserParam,

			// Security
			HasAccess: HasAccess,
		},
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(ds applications.Datasource, prefix string) ([]string, error) {

	linkTables(ds)
	createCache(ds)
	ds.SetModule(assets.MODULEID, assets.VERSION)

	go buildCache(ds)

	return []string{}, nil
}

func Synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	result := []string{}

	ok, res := base.VerifyNeeds(ds, assets.Needs)
	result = append(result, res...)
	if !ok {
		return result, nil
	}

	installed := base.ModuleInstalledVersion(ds, assets.MODULEID)

	// synchro tables
	err, r := synchroTables(ds, installed)
	result = append(result, r...)
	if err != nil {
		return result, err
	}

	cds := ds.CloneShell()
	_, err = cds.StartTransaction()
	if err != nil {
		result = append(result, err.Error())
		return result, err
	}

	// installation or upgrade ?
	if installed != "" {
		err, r = upgrade(cds, installed)
	} else {
		err, r = install(cds)
	}
	result = append(result, r...)
	if err == nil {
		err = base.AddModule(cds, assets.MODULEID, tools.Message(messages, "MODULENAME"), assets.VERSION)
		if err == nil {
			result = append(result, tools.Message(messages, "modulemodified", assets.MODULEID))
			result = append(result, tools.Message(messages, "commit"))
			err = cds.Commit()
			if err != nil {
				result = append(result, err.Error())
			}
			return result, err
		}
	}
	result = append(result, tools.Message(messages, "rollback", err))
	err1 := cds.Rollback()
	if err1 != nil {
		result = append(result, err1.Error())
	}

	return result, err
}

func StartContext(ds applications.Datasource, ctx *context.Context) error {

	sitecontextname, _ := ctx.Sysparams.GetString("sitecontext")
	// if browser module is activated, then ctx.Version has the device.
	// Order preference to seek the device:
	// 1. into ctx.Sessionparams["device"]
	// 2. ctx.Version if sessionparam is not set.
	// NOTE: device is ONLY informative
	// If you use this module, the browser extension for Xamboo should always be activated to set ctx.Version correctly
	device, _ := ctx.Sessionparams.GetString("device")
	if device == "" {
		device = ctx.Version
	}

	VerifyUserSession(ctx, ds, sitecontextname, device)
	return nil
}
