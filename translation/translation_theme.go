package translation

import (
	"github.com/webability-go/xdominion"
)

/*
  TABLE: kl_rutatraduccion:
  All the links and routes of kiwilimon
*/

func kl_traducciontema() *xdominion.XTable {
	t := xdominion.NewXTable("translation_theme", "translation_thm_")

	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // AI, PK

	// name
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// source
	// 1 = tablas de Bdd => tienen entrada en traducciontabla, 2 = PC, 3 = MOB, 4 = GRAPHv5, 5 = Identity Manager
	t.AddField(xdominion.XFieldInteger{Name: "source", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	// link if file or table name
	t.AddField(xdominion.XFieldVarChar{Name: "link", Size: 255})

	// list of fields in table
	t.AddField(xdominion.XFieldVarChar{Name: "fields", Size: 4000})

	return t
}
