package base

import (
	"github.com/webability-go/xdominion"

	serverassets "github.com/webability-go/xamboo/assets"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{
	"base_module",
}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{
	"base_module": baseModule,
}

func linkTables(ds serverassets.Datasource) {

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(ds.GetDatabase())
		ds.SetTable(tbl, table)
	}
}

func synchroTables(ds serverassets.Datasource, oldversion string) (error, []string) {

	result := []string{}

	for _, tbl := range moduletablesorder {

		err, res := SynchroTable(ds, tbl)
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

func install(ds serverassets.Datasource) (error, []string) {

	// do things

	return nil, []string{}
}

func upgrade(ds serverassets.Datasource, oldversion string) (error, []string) {

	if oldversion < "0.0.1" {
		// do things
	}

	return nil, []string{}
}
