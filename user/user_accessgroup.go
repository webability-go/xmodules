package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_accessgroup:
// Groups of accesses
func userAccessGroup() *xdominion.XTable {
	t := xdominion.NewXTable("user_accessgroup", "user_acg_")

	// Key of the accessgroup
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	// Name of the group.
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// description
	t.AddField(xdominion.XFieldText{Name: "description"})

	return t
}
