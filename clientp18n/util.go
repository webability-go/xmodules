package clientp18n

import (
	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xdominion"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{
	//	"wiki_wiki",
}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{
	//	"wiki_wiki":                  wikiwiki,
}

func buildTables(ds applications.Datasource) {

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(ds.GetDatabase())
		ds.SetTable(tbl, table)
	}
}

func createCache(ds applications.Datasource) []string {

	//	ctx.SetCache("wiki:wikis", xcore.NewXCache("wiki:wikis", 0, 0))

	return []string{}
}

func buildCache(ds applications.Datasource) []string {
	/*
		wiki_wiki := ctx.GetTable("wiki_wiki")
		if wiki_wiki == nil {
			return []string{"xmodules::wiki::buildCache: Error, the wiki_wiki table is not available on this context"}
		}
		cache := ctx.GetCache("wiki:wikis")
		if cache == nil {
			return []string{"xmodules::wiki::buildCache: Error, the wiki cache is not available on this site context"}
		}

		// Loads all data in XCache
		wikis, _ := wiki_wiki.SelectAll()

		if wikis != nil {
			for _, m := range *wikis {
				// creates structure on language
				str := CreateStructurewikiByData(ctx, m.Clone())
				key, _ := m.GetString("key")
				cache.Set(key, str)
			}
		}
	*/
	return []string{}
}

func createTables(ds applications.Datasource) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {

		table := ds.GetTable(tbl)
		if table == nil {
			return []string{"xmodules::clientp18n::createTables: Error, the table is not available on this context:" + tbl}
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

func loadTables(ds applications.Datasource) []string {
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
