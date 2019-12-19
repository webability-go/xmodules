package translation

import (
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["translation_theme"] = translationTheme()
	sitecontext.Tables["translation_theme"].SetBase(sitecontext.Databases[databasename])

	sitecontext.Tables["translation_info"] = translationInfo()
	sitecontext.Tables["translation_info"].SetBase(sitecontext.Databases[databasename])
}
