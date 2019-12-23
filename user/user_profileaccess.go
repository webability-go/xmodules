package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_profileaccess:
// All the simple accesses for a profile
func userProfileAccess() *xdominion.XTable {
	t := xdominion.NewXTable("user_profileaccess", "user_pra_")

	// profile
	t.AddField(xdominion.XFieldVarChar{Name: "profile", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_profile", "user_pro_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.DC},
	}})

	// access
	t.AddField(xdominion.XFieldVarChar{Name: "access", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_access", "user_acc_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.DC},
	}})

	return t
}
