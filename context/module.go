package context

import (
	"github.com/webability-go/xcore"
)

func IsModuleAuthorized(sitecontext *Context, id string) bool {
	_, ok := sitecontext.Modules[id]
	return ok
}

// ModuleStatus returns status of the module:
// "" not installed
// "version.number" version installed
func ModuleInstalledVersion(sitecontext *Context, id string) string {
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
