package ingredient

import (
	"github.com/webability-go/xdominion"
)

/*
  Ingredients aisle in supermarket
*/

func ingredientAisle() *xdominion.XTable {
	t := xdominion.NewXTable("ingredient_aisle", "ingredient_ais_")

	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}})

	// name
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 50, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// order
	t.AddField(xdominion.XFieldInteger{Name: "order", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// date last modif
	t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
