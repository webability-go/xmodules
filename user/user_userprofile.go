package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_userprofile:
// All the profiles this user has access
func userUserProfile() *xdominion.XTable {
	t := xdominion.NewXTable("user_userprofile", "user_usp_")

	// user
	t.AddField(xdominion.XFieldInteger{Name: "user", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_user", "user_usr_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.DC},
	}})

	// profile
	t.AddField(xdominion.XFieldVarChar{Name: "profile", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_profile", "user_pro_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.DC},
	}})

	return t
}
