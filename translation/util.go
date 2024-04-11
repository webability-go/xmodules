package translation

import (
	"fmt"

	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/adminmenu"
	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user"
	userassets "github.com/webability-go/xmodules/user/assets"

	"github.com/webability-go/xmodules/translation/assets"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{
	"kl_language",
	"kl_traduccionfuente",
	"kl_traducciontema",
	"kl_traducciontabla",
}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{
	"kl_language":         kl_language,
	"kl_traduccionfuente": kl_traduccionfuente,
	"kl_traducciontema":   kl_traducciontema,
	"kl_traducciontabla":  kl_traducciontabla,
}

func linkTables(ds applications.Datasource) {

	langs := ds.GetLanguages()

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(ds.GetDatabase())
		ds.SetTable(tbl, table)
		table.SetLanguage(langs[0])
	}
}

func createCache(ds applications.Datasource) []string {

	return []string{}
}

func buildCache(ds applications.Datasource) []string {

	return []string{}
}

func synchroTables(ds applications.Datasource, oldversion string) (error, []string) {

	result := []string{}

	for _, tbl := range moduletablesorder {

		err, res := base.SynchroTable(ds, tbl)
		result = append(result, res...)
		if err != nil {
			return err, result
		}
	}

	if oldversion < "0.0.1" {
		// do things
	}

	return nil, result
}

func install(ds applications.Datasource) (error, []string) {

	result := []string{}

	// Accesses
	err := user.AddAccessGroup(ds, &userassets.AccessGroup{
		Key:         assets.ACCESSGROUP,
		Name:        tools.Message(messages, "translationgroup.name"),
		Description: tools.Message(messages, "translationgroup.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         assets.ACCESSLANGUAGE,
		Name:        tools.Message(messages, "translationlanguage.name"),
		Group:       assets.ACCESSGROUP,
		Description: tools.Message(messages, "translationlanguage.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         assets.ACCESSSOURCE,
		Name:        tools.Message(messages, "translationsource.name"),
		Group:       assets.ACCESSGROUP,
		Description: tools.Message(messages, "translationsource.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         assets.ACCESSTHEME,
		Name:        tools.Message(messages, "translationtheme.name"),
		Group:       assets.ACCESSGROUP,
		Description: tools.Message(messages, "translationtheme.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         assets.ACCESS,
		Name:        tools.Message(messages, "translation.name"),
		Group:       assets.ACCESSGROUP,
		Description: tools.Message(messages, "translation.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	err = user.AddAccess(ds, &userassets.Access{
		Key:         assets.ACCESSTOOLS,
		Name:        tools.Message(messages, "translationtools.name"),
		Group:       assets.ACCESSGROUP,
		Description: tools.Message(messages, "translationtools.description"),
	})
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	mainoption := xdominion.XRecord{
		"key":          "_translationadmin",
		"group":        "_admin",
		"father":       nil,
		"name":         tools.Message(messages, "translationfolder.name"),
		"access":       assets.ACCESS,
		"icon1":        "folder.png",
		"call1":        "openclose",
		"description1": tools.Message(messages, "translationfolder.description"),
	}
	err = adminmenu.AddOption(ds, &mainoption)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	option := xdominion.XRecord{
		"key":          "_translationadmin_language",
		"group":        "_admin",
		"father":       "_translationadmin",
		"name":         tools.Message(messages, "translationlanguageoption.name"),
		"access":       assets.ACCESSLANGUAGE,
		"icon1":        "option.png",
		"call1":        "translationadmin/language|single",
		"description1": tools.Message(messages, "translationlanguageoption.description"),
	}
	err = adminmenu.AddOption(ds, &option)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	option = xdominion.XRecord{
		"key":          "_translationadmin_source",
		"group":        "_admin",
		"father":       "_translationadmin",
		"name":         tools.Message(messages, "translationsourceoption.name"),
		"access":       assets.ACCESSSOURCE,
		"icon1":        "option.png",
		"call1":        "translationadmin/source|single",
		"description1": tools.Message(messages, "translationsourceoption.description"),
	}
	err = adminmenu.AddOption(ds, &option)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	option = xdominion.XRecord{
		"key":          "_translationadmin_theme",
		"group":        "_admin",
		"father":       "_translationadmin",
		"name":         tools.Message(messages, "translationthemeoption.name"),
		"access":       assets.ACCESSTHEME,
		"icon1":        "option.png",
		"call1":        "translationadmin/theme|single",
		"description1": tools.Message(messages, "translationthemeoption.description"),
	}
	err = adminmenu.AddOption(ds, &option)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	option = xdominion.XRecord{
		"key":          "_translationadmin_option",
		"group":        "_admin",
		"father":       "_translationadmin",
		"name":         tools.Message(messages, "translationoption.name"),
		"access":       assets.ACCESS,
		"icon1":        "option.png",
		"call1":        "translationadmin|single",
		"description1": tools.Message(messages, "translationoption.description"),
	}
	err = adminmenu.AddOption(ds, &option)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	option = xdominion.XRecord{
		"key":          "_translationadmin_tools",
		"group":        "_admin",
		"father":       "_translationadmin",
		"name":         tools.Message(messages, "translationtoolsoption.name"),
		"access":       assets.ACCESSTOOLS,
		"icon1":        "option.png",
		"call1":        "translationadmin/tools|single",
		"description1": tools.Message(messages, "translationtoolsoption.description"),
	}
	err = adminmenu.AddOption(ds, &option)
	if err != nil {
		result = append(result, err.Error())
		return err, result
	}

	return nil, []string{
		fmt.Sprint("translationadmin options added"),
	}
}

func upgrade(ds applications.Datasource, oldversion string) (error, []string) {

	if oldversion < "0.0.1" {
		// do things
		return install(ds)
	}
	return install(ds)

	return nil, []string{}
}
