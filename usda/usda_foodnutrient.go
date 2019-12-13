package usda

import (
	"github.com/webability-go/xdominion"
)

/*
  Tabla de nutrientes de los alimentos de la USDA
*/

func usda_foodnutrient() *xdominion.XTable {
	t := xdominion.NewXTable("usda_foodnutrient", "usda_fn_")

	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // PK, AI

	t.AddField(xdominion.XFieldVarChar{Name: "food", Size: 5, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"usda_food", "usda_fo_key"}},
	}})

	t.AddField(xdominion.XFieldVarChar{Name: "nutrient", Size: 3, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"usda_nutrient", "usda_nu_key"}},
		xdominion.XConstraint{Type: xdominion.MU, Data: []string{"food"}},
	}})

	t.AddField(xdominion.XFieldFloat{Name: "value", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
