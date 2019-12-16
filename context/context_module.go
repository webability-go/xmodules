package context

import (
	"github.com/webability-go/xdominion"
)

// contextModule contains the list of installed modules
func contextModule() *xdominion.XTable {
	t := xdominion.NewXTable("context_module", "context_mod_")

	// Clave del pais, primary key
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 15, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // AI, PK
	// Nombre
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}}) // NN

	// extra field (3 char ISO field, id of state, number of state, etc) if any
	t.AddField(xdominion.XFieldVarChar{Name: "version", Size: 10})

	return t
}
