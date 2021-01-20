package base

import (
	"errors"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/tools"
)

var ModulesList = &Modules{}

type Modules map[string]assets.Module

func (ml *Modules) Register(m assets.Module) {
	id := m.GetID()
	(*ml)[id] = m
}

func (ml *Modules) Get(id string) assets.Module {
	return (*ml)[id]
}

type Module struct {
	ID      string
	Version string

	Languages map[language.Tag]string
	Needs     []string

	FSetup        func(assets.Datasource, string) ([]string, error)
	FSynchronize  func(assets.Datasource, string) ([]string, error)
	FStartContext func(assets.Datasource, *assets.Context) error
}

func (m *Module) GetID() string {
	return m.ID
}

func (m *Module) GetVersion() string {
	return m.Version
}

func (m *Module) GetLanguages() map[language.Tag]string {
	return m.Languages
}

func (m *Module) GetNeeds() []string {
	return m.Needs
}

func (m *Module) GetInstalledVersion(ds assets.Datasource) string {
	return ModuleInstalledVersion(ds, m.ID)
}

func (m *Module) Setup(ds assets.Datasource, prefix string) ([]string, error) {
	if m.FSetup != nil {
		return m.FSetup(ds, prefix)
	}
	return []string{}, nil
}

func (m *Module) Synchronize(ds assets.Datasource, prefix string) ([]string, error) {
	if m.FSynchronize != nil {
		return m.FSynchronize(ds, prefix)
	}
	return []string{}, nil
}

func (m *Module) StartContext(ds assets.Datasource, ctx *assets.Context) error {
	if m.FStartContext != nil {
		return m.FStartContext(ds, ctx)
	}
	return nil
}

// ModuleStatus returns status of the module:
// "" not installed
// "version.number" version installed
func ModuleInstalledVersion(ds assets.Datasource, id string) string {
	base_module := ds.GetTable("base_module")
	if base_module == nil {
		ds.Log("main", tools.Message(messages, "notable", "base_module"))
		return ""
	}
	data, err := base_module.SelectOne(id)
	if err != nil || data == nil {
		return "" // not installed
	}
	v, _ := data.GetString("version")
	return v
}

func GetModule(ds assets.Datasource, id string) *xcore.XDataset {

	data := xcore.XDataset{}
	data["name"] = ds.GetName()
	data["module"] = id
	data["codeversion"] = ds.GetModule(id)
	data["installedversion"] = ModuleInstalledVersion(ds, id)
	return &data
}

// AddModule will insert a record in the modules table and sends back status error
func AddModule(ds assets.Datasource, id string, name string, version string) error {
	base_module := ds.GetTable("base_module")
	if base_module == nil {
		msgerror := tools.Message(messages, "notable", "base_module")
		ds.Log("main", msgerror)
		return errors.New(msgerror)
	}
	_, err := base_module.Upsert(id, xdominion.XRecord{
		"key":     id,
		"name":    name,
		"version": version,
	}, ds.GetTransaction())
	return err
}
