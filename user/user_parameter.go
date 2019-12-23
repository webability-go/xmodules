package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_profileaccessextended:
// All the extended accesses for a profile
func userParameter() *xdominion.XTable {
	t := xdominion.NewXTable("user_parameter", "user_par_")

	// Key of the parameter, automatic consecutive
	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // PK

	// user
	t.AddField(xdominion.XFieldInteger{Name: "user", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_user", "user_usr_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.DC},
	}})

	// parameter id
	t.AddField(xdominion.XFieldVarChar{Name: "id", Size: 60, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// parameter value
	t.AddField(xdominion.XFieldText{Name: "value", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
