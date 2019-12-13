package translation

import (
	"github.com/webability-go/xdominion"
)

/*
  TABLE: translation_info:
  All the translation words into any language
*/

func translation_info() *xdominion.XTable {
	t := xdominion.NewXTable("trnaslation_info", "translation_info_")

	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // AI, PK

	// theme
	t.AddField(xdominion.XFieldInteger{Name: "theme", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"translation_theme", "translation_thm_key"}},
	}})

	// iso language
	t.AddField(xdominion.XFieldVarChar{Name: "language", Size: 2, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// external key
	t.AddField(xdominion.XFieldVarChar{Name: "externalkey", Size: 50, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// field
	t.AddField(xdominion.XFieldVarChar{Name: "field", Size: 50, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// translation
	t.AddField(xdominion.XFieldVarChar{Name: "translation", Size: 4000, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// date of last modif
	t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// last user
	t.AddField(xdominion.XFieldInteger{Name: "lastuser", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_usuario", "kl_usr_clave"}},
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// verify
	t.AddField(xdominion.XFieldInteger{Name: "verified"})

	// verify user
	t.AddField(xdominion.XFieldInteger{Name: "verifyuser", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_usuario", "kl_usr_clave"}},
	}})

	// verify date
	t.AddField(xdominion.XFieldDateTime{Name: "verifydate"})

	return t
}
