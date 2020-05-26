// Package user contains the list of administrative user for the system.
// All users have accesses, into a profile and even extended access based upon table records.
// It needs base xmodule.
package ingredient

import (
	"golang.org/x/text/language"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID              = "ingredient"
	VERSION               = "0.0.1"
	TRANSLATIONTHEME      = "ingredient"
	TRANSLATIONTHEMEAISLE = "ingaisle"
)

func init() {
	m := &base.Module{
		ID:           MODULEID,
		Version:      VERSION,
		Languages:    map[language.Tag]string{language.English: "Ingredients", language.Spanish: "Ingredientes", language.French: "Ingrédients"},
		Needs:        []string{"base"},
		FSetup:       Setup,
		FSynchronize: Synchronize,
	}
	base.ModulesList.Register(m)
}

// InitModule is called during the init phase to link the module with the system
// adds tables and caches to ctx::database
func Setup(ds serverassets.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	buildTables(ctx)
	createCache(ctx)
	ctx.SetModule(MODULEID, VERSION)

	go buildCache(ctx)

	return []string{}, nil
}

func Synchronize(ds serverassets.Datasource, prefix string) ([]string, error) {

	ctx := ds.(*base.Datasource)
	translation.AddTheme(ctx, TRANSLATIONTHEME, "Ingredients", translation.SOURCETABLE, "", "name,plural")
	translation.AddTheme(ctx, TRANSLATIONTHEMEAISLE, "Ingredients Aisles", translation.SOURCETABLE, "", "name")

	messages := []string{}

	messages = append(messages, "Analysing ingredient_aisle table.")
	num, err := ctx.GetTable("ingredient_aisle").Count(nil)
	if err != nil || num == 0 {
		err1 := ctx.GetTable("ingredient_aisle").Synchronize()
		if err1 != nil {
			messages = append(messages, "The table ingredient_aisle was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table ingredient_aisle was created (again)")
		}
	} else {
		messages = append(messages, "The table ingredient_aisle was not created because it contains data.")
	}

	messages = append(messages, "Analysing ingredient_ingredient table.")
	num, err = ctx.GetTable("ingredient_ingredient").Count(nil)
	if err != nil || num == 0 {
		err1 := ctx.GetTable("ingredient_ingredient").Synchronize()
		if err1 != nil {
			messages = append(messages, "The table ingredient_ingredient was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table ingredient_ingredient was created (again)")
		}
	} else {
		messages = append(messages, "The table ingredient_ingredient was not created because it contains data.")
	}

	// fill metric and translations
	return messages, nil
}

func StartContext(ds serverassets.Datasource, ctx *serverassets.Context) error {
	return nil
}
