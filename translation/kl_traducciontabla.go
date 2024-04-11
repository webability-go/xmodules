package translation

import (
	"github.com/webability-go/xdominion"
)

/*
  TABLE: kl_rutatraduccion:
  All the links and routes of kiwilimon
*/

func kl_traducciontabla() *xdominion.XTable {
	t := xdominion.NewXTable("kl_traducciontabla", "kl_tta_")

	t.AddField(xdominion.XFieldInteger{Name: "clave", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // AI, PK

	// tema
	t.AddField(xdominion.XFieldInteger{Name: "tema", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_traducciontema", "kl_trt_clave"}},
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// idioma
	t.AddField(xdominion.XFieldVarChar{Name: "idioma", Size: 5, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_language", "kl_lan_key"}},
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// clave externa
	t.AddField(xdominion.XFieldVarChar{Name: "claveext", Size: 50, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// campo
	t.AddField(xdominion.XFieldVarChar{Name: "campo", Size: 50, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// link
	t.AddField(xdominion.XFieldVarChar{Name: "traduccion", Size: 4000, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// fecha
	t.AddField(xdominion.XFieldDateTime{Name: "fecha", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// last user
	t.AddField(xdominion.XFieldInteger{Name: "lastuser", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_usuario", "kl_usr_clave"}},
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// verify
	t.AddField(xdominion.XFieldInteger{Name: "verify"})

	// verify user
	t.AddField(xdominion.XFieldInteger{Name: "verifyuser", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_usuario", "kl_usr_clave"}},
	}})

	// verify date
	t.AddField(xdominion.XFieldDateTime{Name: "verifydate"})

	return t
}
