// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package metric

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/metric/assets"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID         = "metric"
	VERSION          = "0.0.1"
	TRANSLATIONTHEME = "metric"
)

var ModuleMetric = assets.ModuleEntries{}

func init() {
	m := &base.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Metrics", language.Spanish: "Métricas", language.French: "Métriques"},
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
	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ds applications.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	translation.AddTheme(ctx, TRANSLATIONTHEME, "Metric units", translation.SOURCETABLE, "", "name,plural")

	messages := []string{}
	messages = append(messages, "Analysing metric_unit table.")
	num, err := ctx.GetTable("metric_unit").Count(nil)
	if err != nil || num == 0 {
		err1 := ctx.GetTable("metric_unit").Synchronize()
		if err1 != nil {
			messages = append(messages, "The table metric_unit was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table metric_unit was created (again)")
		}
	} else {
		messages = append(messages, "The table metric_unit was not created because it contains data.")
	}

	// fill metric and translations
	return messages, nil
}

func StartContext(ds applications.Datasource, ctx *context.Context) error {
	return nil
}
