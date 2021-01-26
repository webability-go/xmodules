package useradmin

import (
	serverassets "github.com/webability-go/xamboo/assets"

	//	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/adminmenu"
)

func install(ds serverassets.Datasource) (error, []string) {

	// Group, just in case (upsert)
	adminmenu.AddGroup(ds, "_admin", "System Administration")

	mainoption := xdominion.XRecord{
		"key":          "user",
		"group":        "admin",
		"father":       nil,
		"name":         "User administration",
		"access":       nil,
		"icon1":        nil,
		"call1":        "openclose",
		"description1": "Open or close the option",
	}
	err := adminmenu.AddOption(ds, &mainoption)

	if err != nil {
		ds.Log("main", "Error inserting admin adminmenu_option", err)
		return err, []string{"xmodules::adminmenu::loadTables: Error upserting the admin adminmenu group"}
	}

	option := xdominion.XRecord{
		"key":          "userlist",
		"group":        "admin",
		"father":       "user",
		"name":         "User list",
		"access":       nil,
		"icon1":        nil,
		"call1":        "openclose",
		"description1": "Open or close the option",
	}
	err = adminmenu.AddOption(ds, &option)

	if err != nil {
		ds.Log("main", "Error inserting admin adminmenu_option", err)
		return err, []string{"xmodules::adminmenu::loadTables: Error upserting the admin adminmenu group"}
	}

	return nil, []string{"ok"}
}

func upgrade(ds serverassets.Datasource, oldversion string) (error, []string) {

	if oldversion < "0.0.1" {
		// do things
	}

	return nil, []string{}
}
