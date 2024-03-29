// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package translation

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation/assets"
)

var ModuleEntries = assets.ModuleEntries{
	TranslateOne: TranslateOne,
}

func init() {
	m := &base.Module{
		ID:           assets.MODULEID,
		Version:      assets.VERSION,
		Languages:    map[language.Tag]string{language.English: "Translation tables", language.Spanish: "Tablas de traducción", language.French: "Tables de traduction"},
		Needs:        []string{"base"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to ctx::database
func Setup(ds applications.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	buildTables(ctx)
	ctx.SetModule(assets.MODULEID, assets.VERSION)

	return []string{}, nil
}

func Synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	messages := []string{}

	ctx := ds.(*base.Datasource)
	// Needed modules: base
	vc := base.ModuleInstalledVersion(ctx, "base")
	if vc == "" {
		messages = append(messages, "xmodules/base need to be installed before installing xmodules/translation.")
		return messages, nil
	}
	vc1 := base.ModuleInstalledVersion(ctx, "user")
	vc2 := base.ModuleInstalledVersion(ctx, "userlink")
	if vc1 == "" && vc2 == "" {
		messages = append(messages, "xmodules/user or xmodules/userlink need to be installed before installing xmodules/translation.")
		return messages, nil
	}

	// create tables
	messages = append(messages, createTables(ctx)...)

	// Inserting into base-modules
	// Be sure base module is on db: fill base module (we should get this from xmodule.conf)
	err := base.AddModule(ctx, assets.MODULEID, "Multilanguages translation tables for Xamboo", assets.VERSION)
	if err == nil {
		messages = append(messages, "The entry "+assets.MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+assets.MODULEID+" in the modules table: "+err.Error())
	}

	return messages, nil
}

func StartContext(ds applications.Datasource, ctx *context.Context) error {
	return nil
}

func TranslateOne(input string) (string, error) {
	return "TRANSLATED: {" + input + "}", nil
}
