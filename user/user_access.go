package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_access:
// All the simple accesses for a user
func userAccess() *xdominion.XTable {
	t := xdominion.NewXTable("user_access", "user_acc_")

	// Key of the accessgroup
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	// Name of the group.
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// access group, optional
	t.AddField(xdominion.XFieldVarChar{Name: "group", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_accessgroup", "user_acg_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// description
	t.AddField(xdominion.XFieldText{Name: "description"})

	return t
}
