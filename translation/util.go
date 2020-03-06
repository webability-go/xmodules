package translation

import (
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{
	"translation_theme",
	"translation_info",
}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{
	"translation_theme": translationTheme,
	"translation_info":  translationInfo,
}

func buildTables(sitecontext *context.Context, databasename string) {

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(sitecontext.GetDatabase(databasename))
		sitecontext.SetTable(tbl, table)
	}
}

func createTables(sitecontext *context.Context) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {

		table := sitecontext.GetTable(tbl)
		if table == nil {
			return []string{"xmodules::translation::createTables: Error, the table is not available on this context:" + tbl}
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
