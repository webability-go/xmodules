package translation

import (
	"github.com/webability-go/xdominion"
)

// translationTheme table: All the translation group of words
func translationTheme() *xdominion.XTable {
	t := xdominion.NewXTable("translation_theme", "translation_thm_")

	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 10, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}})

	// name
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// source
	// 1 = tablas de Bdd => tienen entrada en traducciontabla, 2 = file
	t.AddField(xdominion.XFieldInteger{Name: "source", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	// link if file or table name
	t.AddField(xdominion.XFieldVarChar{Name: "link", Size: 255})

	// list of fields in table
	t.AddField(xdominion.XFieldVarChar{Name: "fields", Size: 4000})

	return t
}
