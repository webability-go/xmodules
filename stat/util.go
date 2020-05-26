package stat

import (
	"strconv"
	"time"

	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/base"
)

func buildTables(ctx *base.Datasource, prefix string) {

	// open 12 tables for each file
	for i := 1; i < 13; i++ {
		m := ""
		if i < 10 {
			m = "0"
		}
		m += strconv.Itoa(i)

		table := stat_stat(prefix, m)
		table.SetBase(ctx.GetDatabase())
		table.SetLanguage(language.English)
		ctx.SetTable(prefix+"stat_"+m, table)
	}
}

func getMonth() string {
	currentTime := time.Now()
	return currentTime.Format("01")
}
