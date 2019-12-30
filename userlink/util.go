package userlink

import (
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{
	"user_user",
}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{
	"user_user": userUser,
}

func buildTables(sitecontext *context.Context, databasename string) {

	for _, tbl := range moduletablesorder {
		sitecontext.Tables[tbl] = moduletables[tbl]()
		sitecontext.Tables[tbl].SetBase(sitecontext.Databases[databasename])
	}
}

func createTables(sitecontext *context.Context) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {
		messages = append(messages, "Analysing "+tbl+" table.")
		num, err := sitecontext.Tables[tbl].Count(nil)
		if err != nil || num == 0 {
			err1 := sitecontext.Tables[tbl].Synchronize()
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
