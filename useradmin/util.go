package useradmin

import (
	serverassets "github.com/webability-go/xamboo/assets"

	//	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/adminmenu"
)

func loadTables(ds serverassets.Datasource) ([]string, error) {

	// Insert menu
	adminmenu.AddGroup(ds, "admin", "System Administration")

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
		return []string{"xmodules::adminmenu::loadTables: Error upserting the admin adminmenu group"}, err
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
		return []string{"xmodules::adminmenu::loadTables: Error upserting the admin adminmenu group"}, err
	}

	return []string{"ok"}, nil
}
