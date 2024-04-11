// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package chatgpt

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"
	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"

	"github.com/webability-go/xmodules/chatgpt/assets"
)

//go:embed languages/*.language
var fsmessages embed.FS
var messages *map[language.Tag]*xcore.XLanguage

//go:embed pages static
var files embed.FS

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
		Entries: &assets.ModuleEntries{
			Translate: Translate,
			Ask:       Ask,
		},
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to ctx::database
func Setup(ds applications.Datasource, prefix string) ([]string, error) {

	linkTables(ds)
	createCache(ds)
	ds.SetModule(assets.MODULEID, assets.VERSION)

	return []string{}, nil
}

func Synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	result := []string{}

	ok, res := base.VerifyNeeds(ds, assets.Needs)
	result = append(result, res...)
	if !ok {
		return result, nil
	}

	installed := base.ModuleInstalledVersion(ds, assets.MODULEID)

	// synchro tables
	err, r := synchroTables(ds, installed)
	result = append(result, r...)
	if err != nil {
		return result, err
	}

	cds := ds.CloneShell()
	_, err = cds.StartTransaction()
	if err != nil {
		result = append(result, err.Error())
		return result, err
	}

	// installation or upgrade ?
	if installed != "" {
		err, r = upgrade(cds, installed)
	} else {
		err, r = install(cds)
	}
	result = append(result, r...)
	if err == nil {
		err = base.AddModule(cds, assets.MODULEID, tools.Message(messages, "MODULENAME"), assets.VERSION)
		if err == nil {
			result = append(result, tools.Message(messages, "modulemodified", assets.MODULEID))
			result = append(result, tools.Message(messages, "commit"))
			err = cds.Commit()
			if err != nil {
				result = append(result, err.Error())
			}
		} else {
			result = append(result, tools.Message(messages, "rollback", err))
			err = cds.Rollback()
			if err != nil {
				result = append(result, err.Error())
			}
		}
	}
	/*
		// copy files
		pathadmin, _ := cds.GetConfig().GetString("pathinstalladmin")
		pages, _ := fs.Sub(files, "pages")
		err, rsf := base.SynchroFiles(pages, pathadmin)
		result = append(result, rsf...)
		if err != nil {
			result = append(result, err.Error())
		}

		pathadminstatic, _ := cds.GetConfig().GetString("pathinstalladminstatic")
		static, _ := fs.Sub(files, "static")
		err, rsf = base.SynchroFiles(static, pathadminstatic)
		result = append(result, rsf...)
		if err != nil {
			result = append(result, err.Error())
		}
	*/
	return result, nil
}

func StartContext(ds applications.Datasource, ctx *context.Context) error {
	return nil
}
