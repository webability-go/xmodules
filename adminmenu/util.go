package adminmenu

import (
	"fmt"
	//	"time"

	//	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/base"
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

func buildTables(ds *base.Datasource) {

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(ds.GetDatabase())
		ds.SetTable(tbl, table)
	}
}

func createCache(ds *base.Datasource) []string {

	return []string{}
}

func buildCache(ds *base.Datasource) []string {

	return []string{}
}

func createTables(ds *base.Datasource) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {

		table := ds.GetTable(tbl)
		if table == nil {
			return []string{"xmodules::adminmenu::createTables: Error, the table is not available on this datasource:" + tbl}
		}

		messages = append(messages, "Analysing "+tbl+" table.")
		num, err := table.Count(nil)
		if err != nil || num == 0 {
			err1 := table.Synchronize()
			if err1 != nil {
				messages = append(messages, "The table "+tbl+" was not created: "+err1.Error())
			} else {
				messages = append(messages, "The table "+tbl+" was created (again)")
			}
		} else {
			messages = append(messages, "The table "+tbl+" was not created because it contains data.")
		}
	}

	return messages
}

func loadTables(ds *base.Datasource) []string {

	adminmenu_group := ds.GetTable("adminmenu_group")
	if adminmenu_group == nil {
		return []string{"xmodules::adminmenu::createTables: Error, the table adminmenu_group is not available on this datasource"}
	}

	adminmenu_option := ds.GetTable("adminmenu_option")
	if adminmenu_option == nil {
		return []string{"xmodules::adminmenu::createTables: Error, the table adminmenu_option is not available on this datasource"}
	}

	// insert admin group
	_, err := adminmenu_group.Upsert("admin", xdominion.XRecord{
		"key":  "admin",
		"name": "Administration menu",
	})
	if err != nil {
		ds.Log("main", "Error inserting admin adminmenu_group", err)
		return []string{"xmodules::adminmenu::loadTables: Error upserting the admin adminmenu group"}
	}

	return []string{
		fmt.Sprint(adminmenu_option.Count(nil)) + " admin menu options added",
	}
}
