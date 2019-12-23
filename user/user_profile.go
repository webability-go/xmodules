package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_profile:
// The differents profiles for administrators
// The profiles are cumulatives for administrators
func userProfile() *xdominion.XTable {
	t := xdominion.NewXTable("user_profile", "user_pro_")

	// Key of the profile
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	// The status of the profile. Any profile can be desactivated at any moment.
	// 1 = ok, 2 = desactivated
	t.AddField(xdominion.XFieldInteger{Name: "status", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// Name of the profile.
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// description
	t.AddField(xdominion.XFieldText{Name: "description"})

	return t
}
