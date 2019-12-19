package context

import (
	"github.com/webability-go/xcore"
	"github.com/webability-go/xdominion"
)

func IsModuleAuthorized(sitecontext *Context, id string) bool {
	_, ok := sitecontext.Modules[id]
	return ok
}

// ModuleStatus returns status of the module:
// "" not installed
// "version.number" version installed
func ModuleInstalledVersion(sitecontext *Context, id string) string {
	if sitecontext.Tables["context_module"] == nil {
		return ""
	}
	data, err := sitecontext.Tables["context_module"].SelectOne(id)
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
	data["codeversion"] = sitecontext.Modules[id]
	data["installedversion"] = ModuleInstalledVersion(sitecontext, id)
	return &data
}

// AddModule will insert a record in the modules table and sends back status error
func AddModule(sitecontext *Context, id string, name string, version string) error {
	_, err := sitecontext.Tables["context_module"].Upsert(id, xdominion.XRecord{
		"key":     id,
		"name":    name,
		"version": version,
	})
	return err
}
