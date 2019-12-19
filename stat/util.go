package stat

import (
	"strconv"
	"time"

	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, prefix string, databasename string) {

	// open 12 tables for each file
	for i := 1; i < 13; i++ {
		m := ""
		if i < 10 {
			m = "0"
		}
		m += strconv.Itoa(i)

		sitecontext.Tables[prefix+"stat_"+m] = stat_stat(prefix, m)
		sitecontext.Tables[prefix+"stat_"+m].SetBase(sitecontext.Databases[databasename])
		sitecontext.Tables[prefix+"stat_"+m].SetLanguage(language.English)
	}
}

func getMonth() string {
	currentTime := time.Now()
	return currentTime.Format("01")
}
