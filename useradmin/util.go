package useradmin

import (
	"fmt"

	"github.com/webability-go/xamboo/applications"

	//	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/adminmenu"
	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user"
	userassets "github.com/webability-go/xmodules/user/assets"
)

func install(ds applications.Datasource) (error, []string) {

	result := []string{}

	// Accesses
	err := user.AddAccessGroup(ds, &userassets.AccessGroup{
		Key:         "_useradmin",
		Name:        tools.Message(messages, "accessgroup.name"),
		Description: tools.Message(messages, "accessgroup.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         "_useradmin",
		Name:        tools.Message(messages, "mainuser.name"),
		Group:       "_useradmin",
		Description: tools.Message(messages, "mainuser.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         "_useradminaccess",
		Name:        tools.Message(messages, "access.name"),
		Group:       "_useradmin",
		Description: tools.Message(messages, "access.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         "_useradminprofile",
		Name:        tools.Message(messages, "profile.name"),
		Group:       "_useradmin",
		Description: tools.Message(messages, "profile.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         "_useradminuser",
		Name:        tools.Message(messages, "user.name"),
		Group:       "_useradmin",
		Description: tools.Message(messages, "user.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	mainoption := xdominion.XRecord{
		"key":          "_useradmin",
		"group":        "_admin",
		"father":       nil,
		"name":         tools.Message(messages, "userfolder.name"),
		"access":       "_useradmin",
		"icon1":        "folder.png",
		"call1":        "openclose",
		"description1": tools.Message(messages, "userfolder.description"),
	}
	err = adminmenu.AddOption(ds, &mainoption)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	option := xdominion.XRecord{
		"key":          "_useradmin_access",
		"group":        "_admin",
		"father":       "_useradmin",
		"name":         tools.Message(messages, "useraccess.name"),
		"access":       "_useradminaccess",
		"icon1":        "option.png",
		"call1":        "useradmin/access|single",
		"description1": tools.Message(messages, "useraccess.description"),
	}
	err = adminmenu.AddOption(ds, &option)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	option = xdominion.XRecord{
		"key":          "_useradmin_profile",
		"group":        "_admin",
		"father":       "_useradmin",
		"name":         tools.Message(messages, "userprofile.name"),
		"access":       "_useradminprofile",
		"icon1":        "option.png",
		"call1":        "useradmin/profile|single",
		"description1": tools.Message(messages, "userprofile.description"),
	}
	err = adminmenu.AddOption(ds, &option)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	option = xdominion.XRecord{
		"key":          "_useradmin_user",
		"group":        "_admin",
		"father":       "_useradmin",
		"name":         tools.Message(messages, "useruser.name"),
		"access":       "_useradminuser",
		"icon1":        "option.png",
		"call1":        "useradmin/user|single",
		"description1": tools.Message(messages, "useruser.description"),
	}
	err = adminmenu.AddOption(ds, &option)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	return nil, []string{
		fmt.Sprint("useradmin options added"),
	}
}

func upgrade(ds applications.Datasource, oldversion string) (error, []string) {

	if oldversion <= "0.0.1" {
		// do things
		return install(ds)
	}

	return nil, []string{}
}
