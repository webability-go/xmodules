package translation

import (
	"github.com/webability-go/xdominion"
)

/*
  TABLE: kl_rutatraduccion:
  All the links and routes of kiwilimon
*/

func kl_traducciontema() *xdominion.XTable {
	t := xdominion.NewXTable("kl_traducciontema", "kl_trt_")

	t.AddField(xdominion.XFieldInteger{Name: "clave", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // AI, PK

	// nombre
	t.AddField(xdominion.XFieldVarChar{Name: "nombre", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// fuente
	// 1 = tablas de Bdd => tienen entrada en traducciontabla, 2 = language Files, 3 = specific files (JS)
	t.AddField(xdominion.XFieldInteger{Name: "fuente", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_traduccionfuente", "kl_trf_clave"}},
	}})

	// config
	t.AddField(xdominion.XFieldText{Name: "config"})

	return t
}
