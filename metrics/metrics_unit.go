package metrics

import (
	"github.com/webability-go/xdominion"
)

/*
  ALl purpose units
*/

func metric_unit() *xdominion.XTable {
	t := xdominion.NewXTable("metric_unit", "metric_unt_")

	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}})

	// name of unit
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// plural
	t.AddField(xdominion.XFieldVarChar{Name: "plural", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// 1 = IS (International System), 2 = US system, 3 = cooking system, 4 = insignificant, 5 = not quantificable
	t.AddField(xdominion.XFieldInteger{Name: "type"})

	// 1 if it is an official metric (meter, liter, gram, gallon, pound, cup... )
	t.AddField(xdominion.XFieldInteger{Name: "oficial"})

	// IS unit for conversion
	t.AddField(xdominion.XFieldInteger{Name: "isunit", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"metric_unit", "metric_unt_key"}},
	}})

	// factor for conversion to SI unit
	t.AddField(xdominion.XFieldFloat{Name: "factorconversion"})

	return t
}
