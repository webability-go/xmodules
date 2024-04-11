package stat

import (
	"strconv"
	"time"

	"github.com/webability-go/xamboo/applications"
	"golang.org/x/text/language"
)

func buildTables(ds applications.Datasource, prefix string) {

	// open 12 tables for each file
	for i := 1; i < 13; i++ {
		m := ""
		if i < 10 {
			m = "0"
		}
		m += strconv.Itoa(i)

		table := stat_stat(prefix, m)
		table.SetBase(ds.GetDatabase())
		table.SetLanguage(language.English)
		ds.SetTable(prefix+"stat_"+m, table)
	}
}

func getMonth() string {
	currentTime := time.Now()
	return currentTime.Format("01")
}
