package adminmenu

import (
	"fmt"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user"
	userassets "github.com/webability-go/xmodules/user/assets"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{
	"adminmenu_group",
	"adminmenu_option",
}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{
	"adminmenu_group":  adminmenuGroup,
	"adminmenu_option": adminmenuOption,
}

func linkTables(ds applications.Datasource) {

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(ds.GetDatabase())
		ds.SetTable(tbl, table)
	}
}

func synchroTables(ds applications.Datasource, oldversion string) (error, []string) {

	result := []string{}

	for _, tbl := range moduletablesorder {

		err, res := base.SynchroTable(ds, tbl)
		result = append(result, res...)
		if err != nil {
			return err, result
		}
	}

	if oldversion < "0.0.1" {
		// do things
	}

	return nil, result
}

func install(ds applications.Datasource) (error, []string) {

	result := []string{}

	// Group, just in case (upsert)
	err := AddGroup(ds, "_admin", tools.Message(messages, "MAINMENU"))
	if err != nil {
		return err, result
	}
	// Accesses
	err = user.AddAccessGroup(ds, &userassets.AccessGroup{
		Key:         "_adminmenu",
		Name:        tools.Message(messages, "accessgroup.name"),
		Description: tools.Message(messages, "accessgroup.description"),
	})
	if err != nil {
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         "_adminmenu",
		Name:        tools.Message(messages, "access.name"),
		Group:       "_adminmenu",
		Description: tools.Message(messages, "access.description"),
	})
	if err != nil {
		return err, result
	}

	mainoption := xdominion.XRecord{
		"key":          "_adminmenu",
		"group":        "_admin",
		"father":       nil,
		"name":         tools.Message(messages, "menufolder.name"),
		"access":       "_adminmenu",
		"icon1":        "folder.png",
		"call1":        "openclose",
		"description1": tools.Message(messages, "menufolder.description1"),
	}
	err = AddOption(ds, &mainoption)
	if err != nil {
		return err, result
	}

	option := xdominion.XRecord{
		"key":          "_adminmenu_group",
		"group":        "_admin",
		"father":       "_adminmenu",
		"name":         tools.Message(messages, "adminmenugroup.name"),
		"access":       "_adminmenu",
		"icon1":        "option.png",
		"call1":        "adminmenu/group",
		"description1": tools.Message(messages, "adminmenugroup.description1"),
	}
	err = AddOption(ds, &option)
	if err != nil {
		return err, result
	}

	option = xdominion.XRecord{
		"key":          "_adminmenu_option",
		"group":        "_admin",
		"father":       "_adminmenu",
		"name":         tools.Message(messages, "adminmenuoption.name"),
		"access":       "_adminmenu",
		"icon1":        "option.png",
		"call1":        "adminmenu/option",
		"description1": tools.Message(messages, "adminmenuoption.description1"),
	}
	err = AddOption(ds, &option)
	if err != nil {
		return err, result
	}

	return nil, []string{
		fmt.Sprint(" admin menu options added"),
	}
}

func upgrade(ds applications.Datasource, oldversion string) (error, []string) {

	if oldversion < "0.0.1" {
		// do things
	}

	return nil, []string{}
}
