// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package usda

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"

	"github.com/webability-go/xmodules/usda/assets"
)

//go:embed languages/*.language
var fsmessages embed.FS
var messages *map[language.Tag]*xcore.XLanguage

func init() {
	messages = tools.BuildMessagesFS(fsmessages, "languages")
	m := &base.Module{
		ID:      assets.MODULEID,
		Version: assets.VERSION,
		Languages: map[language.Tag]string{
			language.English: tools.Message(messages, "MODULENAME", language.English),
			language.Spanish: tools.Message(messages, "MODULENAME", language.Spanish),
		},
		Needs:         assets.Needs,
		FSetup:        Setup,
		FSynchronize:  Synchronize,
		FStartContext: StartContext,
		Entries:       &assets.ModuleEntries{},
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to sitecontext::database
func Setup(ds applications.Datasource, prefix string) ([]string, error) {

	buildTables(ds)
	createCache(ds)
	ds.SetModule(assets.MODULEID, assets.VERSION)

	go buildCache(ds)

	return []string{}, nil
}

func Synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	//	translation.AddTheme(ds, TRANSLATIONTHEME, "USDA nutrients", translation.SOURCETABLE, "", "name,tag")

	messages := []string{}
	messages = append(messages, createTables(ds)...)

	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(ds, assets.MODULEID, "List of USDA food and nutrients", assets.VERSION)
	if err == nil {
		messages = append(messages, "The entry "+assets.MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+assets.MODULEID+" in the modules table: "+err.Error())
	}

	messages = append(messages, loadTables(ds, prefix)...)
	messages = append(messages, buildCache(ds)...)

	return messages, nil
}

func StartContext(ds applications.Datasource, ctx *context.Context) error {
	return nil
}
