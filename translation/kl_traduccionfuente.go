package translation

import (
	"github.com/webability-go/xdominion"
)

/*
  TABLE: kl_traduccionfuente:
  All sources of translations
*/

func kl_traduccionfuente() *xdominion.XTable {
	t := xdominion.NewXTable("kl_traduccionfuente", "kl_trf_")

	t.AddField(xdominion.XFieldInteger{Name: "clave", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // AI, PK

	// nombre
	t.AddField(xdominion.XFieldVarChar{Name: "nombre", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// config
	t.AddField(xdominion.XFieldText{Name: "config"})

	return t
}
