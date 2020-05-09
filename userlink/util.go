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

func buildTables(sitecontext *context.Context) {

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(sitecontext.GetDatabase())
		sitecontext.SetTable(tbl, table)
	}
}

func createTables(sitecontext *context.Context) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {

		table := sitecontext.GetTable(tbl)
		if table == nil {
			return []string{"xmodules::userlink::createTables: Error, the table is not available on this context:" + tbl}
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

func loadTables(ctx *context.Context) []string {
	/*
		wiki_wiki := ctx.GetTable("wiki_wiki")
		if wiki_wiki == nil {
			return []string{"xmodules::wiki::createTables: Error, the table wiki_wiki is not available on this context"}
		}

		if err != nil {
			ctx.Log("main", "Error inserting admin wiki", err)
			return []string{"xmodules::wiki::loadTables: Error upserting the admin wiki"}
		}
	*/
	return []string{}
}
