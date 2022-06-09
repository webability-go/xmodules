package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_profileaccess:
// All the accesses for a direct user
func userUserAccess() *xdominion.XTable {
	t := xdominion.NewXTable("user_useraccess", "user_usa_")

	// user
	t.AddField(xdominion.XFieldInteger{Name: "user", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_user", "user_usr_key"}},
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

	// 1 if this access is denied, or 0
	t.AddField(xdominion.XFieldInteger{Name: "denied", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
