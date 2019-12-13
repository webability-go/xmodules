package usda

import (
	"github.com/webability-go/xdominion"
)

/*
  Tabla de grupos de alimentos de la USDA
*/

func usda_group() *xdominion.XTable {
	t := xdominion.NewXTable("usda_group", "usda_gr_")

	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 4, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 60, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
