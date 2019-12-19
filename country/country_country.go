package country

import (
	"github.com/webability-go/xdominion"
)

/*
  Ingredients for recipes
*/

func countryCountry() *xdominion.XTable {
	t := xdominion.NewXTable("country_country", "country_cou_")

	// Clave del pais, primary key
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 8, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // AI, PK
	// Nombre
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}}) // NN

	// type: 1 = region, 2 = subregion, 3 = country, 4 = subdivision
	t.AddField(xdominion.XFieldInteger{Name: "type", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// padre
	t.AddField(xdominion.XFieldVarChar{Name: "father", Size: 5, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"country_country", "country_cou_key"}},
	}}) // FK: kl_pais

	// extra field (3 char ISO field, id of state, number of state, etc) if any
	t.AddField(xdominion.XFieldVarChar{Name: "iso3", Size: 3})

	// country code
	t.AddField(xdominion.XFieldInteger{Name: "code"})

	// php geo ip extra id if any (states)
	t.AddField(xdominion.XFieldVarChar{Name: "geoip", Size: 3})

	return t
}
