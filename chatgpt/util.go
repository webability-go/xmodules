package chatgpt

import (
	"fmt"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/base"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{}

func linkTables(ds applications.Datasource) {

	langs := ds.GetLanguages()

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(ds.GetDatabase())
		ds.SetTable(tbl, table)
		table.SetLanguage(langs[0])
	}
}

func createCache(ds applications.Datasource) []string {

	return []string{}
}

func buildCache(ds applications.Datasource) []string {

	return []string{}
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

	//	result := []string{}

	/*
		// Accesses
		err := user.AddAccessGroup(ds, &userassets.AccessGroup{
			Key:         assets.ACCESSGROUP,
			Name:        tools.Message(messages, "translationgroup.name"),
			Description: tools.Message(messages, "translationgroup.description"),
		})
		if err != nil {
			result = append(result, err.Error())
			return err, result
		}

		err = user.AddAccess(ds, &userassets.Access{
			Key:         assets.ACCESS,
			Name:        tools.Message(messages, "translationlanguage.name"),
			Group:       assets.ACCESSGROUP,
			Description: tools.Message(messages, "translationlanguage.description"),
		})
		if err != nil {
			result = append(result, err.Error())
			return err, result
		}
	*/

	return nil, []string{
		fmt.Sprint("chatgpt options added"),
	}
}

func upgrade(ds applications.Datasource, oldversion string) (error, []string) {

	if oldversion < "0.0.1" {
		// do things
		return install(ds)
	}
	return install(ds)
}
