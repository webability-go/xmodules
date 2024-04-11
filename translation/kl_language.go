package translation

import (
	"github.com/webability-go/xdominion"
)

/*
  TABLE: kl_rutatraduccion:
  All the links and routes of kiwilimon
*/

func kl_language() *xdominion.XTable {
	t := xdominion.NewXTable("kl_language", "kl_lan_")

	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 5, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // AI, PK

	// Name of the language
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// Is the default language 1 = yes, 0 = no
	t.AddField(xdominion.XFieldInteger{Name: "default", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// config
	t.AddField(xdominion.XFieldText{Name: "config"})

	return t
}
