package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_profileaccessextended:
// All the extended accesses for a profile
func userUserAccessExtended() *xdominion.XTable {
	t := xdominion.NewXTable("user_useraccessextended", "user_uae_")

	// user
	t.AddField(xdominion.XFieldInteger{Name: "user", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_user", "user_usr_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.DC},
	}})

	// access extended
	t.AddField(xdominion.XFieldVarChar{Name: "accessextended", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_accessextended", "user_ace_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.DC},
	}})

	// external key
	t.AddField(xdominion.XFieldVarChar{Name: "keyextended", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// 1 if this access is denied
	t.AddField(xdominion.XFieldInteger{Name: "denied", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
