package ingredient

import (
	"github.com/webability-go/xdominion"
)

/*
  Ingredients for recipes
*/

func ingredientIngredient() *xdominion.XTable {
	t := xdominion.NewXTable("ingredient_ingredient", "ingrediente_ing_")

	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // AI, PK

	// nombre
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// plural
	t.AddField(xdominion.XFieldVarChar{Name: "plural", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// pasillo
	t.AddField(xdominion.XFieldInteger{Name: "aisle", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"ingredient_aisle", "ingredient_ais_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// padre
	t.AddField(xdominion.XFieldInteger{Name: "father", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"ingredient_ingredient", "ingrediente_ing_key"}},
	}})

	// usda
	t.AddField(xdominion.XFieldInteger{Name: "usda", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"usda_food", "usda_fo_key"}},
	}})

	// Unidad SI para pieza
	t.AddField(xdominion.XFieldInteger{Name: "isunit", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"metric_unit", "metric_unt_key"}},
	}})

	// Cantidad en unidad SI para pieza
	t.AddField(xdominion.XFieldFloat{Name: "quantity"})

	// densidad si volumen
	t.AddField(xdominion.XFieldFloat{Name: "density"})

	// date last modif
	t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})
	return t
}
