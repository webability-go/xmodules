package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_accessextended:
// Extended access levels for user for security module
// The extended levels are used for accesses applied on table records.
func userAccessExtended() *xdominion.XTable {
	t := xdominion.NewXTable("user_accessextended", "user_ace_")

	// Key of the access
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	// Name of the access.
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// access group, optional
	t.AddField(xdominion.XFieldVarChar{Name: "group", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_accessgroup", "user_acg_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	// queries
	t.AddField(xdominion.XFieldText{Name: "queries", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// description
	t.AddField(xdominion.XFieldText{Name: "description"})

	return t
}
