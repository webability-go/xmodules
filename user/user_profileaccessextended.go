package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_profileaccessextended:
// All the extended accesses for a profile
func userProfileAccessExtended() *xdominion.XTable {
	t := xdominion.NewXTable("user_profileaccessextended", "user_pae_")

	// profile
	t.AddField(xdominion.XFieldVarChar{Name: "profile", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_profile", "user_pro_key"}},
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

	return t
}
