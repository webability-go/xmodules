package context

import (
	"errors"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
)

func IsModuleAuthorized(sitecontext *Context, id string) bool {
	return sitecontext.GetModule(id) != ""
}

// ModuleStatus returns status of the module:
// "" not installed
// "version.number" version installed
func ModuleInstalledVersion(sitecontext *Context, id string) string {
	context_module := sitecontext.GetTable("context_module")
	if context_module == nil {
		sitecontext.Log("main", "Error: the context_module table is not available within the context xmodule")
		return ""
	}
	data, err := context_module.SelectOne(id)
	if err != nil || data == nil {
		return "" // not installed
	}
	v, _ := data.GetString("version")
	return v
}

func GetModule(sitecontext *Context, id string) *xcore.XDataset {

	data := xcore.XDataset{}
	data["context"] = sitecontext.Name
	data["module"] = id
	data["codeversion"] = sitecontext.GetModule(id)
	data["installedversion"] = ModuleInstalledVersion(sitecontext, id)
	return &data
}

// AddModule will insert a record in the modules table and sends back status error
func AddModule(sitecontext *Context, id string, name string, version string) error {
	context_module := sitecontext.GetTable("context_module")
	if context_module == nil {
		return errors.New("Error: the context_module table is not available within the context xmodule")
	}
	_, err := context_module.Upsert(id, xdominion.XRecord{
		"key":     id,
		"name":    name,
		"version": version,
	})
	return err
}
