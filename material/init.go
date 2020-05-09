// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package material

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID         = "material"
	VERSION          = "0.0.1"
	TRANSLATIONTHEME = "material"
)

func init() {
	m := &context.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Material", language.Spanish: "Material", language.French: "Materiel"},
		Needs:        []string{"context"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	context.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to ctx::database
func Setup(ctx *context.Context, prefix string) ([]string, error) {

	buildTables(ctx)
	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ctx *context.Context, prefix string) ([]string, error) {

	translation.AddTheme(ctx, TRANSLATIONTHEME, "Material", translation.SOURCETABLE, "", "name,plural")

	messages := []string{}
	messages = append(messages, "Analysing material_material table.")
	num, err := ctx.GetTable("material_material").Count(nil)
	if err != nil || num == 0 {
		err1 := ctx.GetTable("material_material").Synchronize()
		if err1 != nil {
			messages = append(messages, "The table material_material was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table material_material was created (again)")
		}
	} else {
		messages = append(messages, "The table material_material was not created because it contains data.")
	}

	// fill metric and translations
	return messages, nil
}