package adminmenu

import (
	"github.com/webability-go/xdominion"
)

// TABLE: adminmenu_option:
// The menu tree for administration
func adminmenuGroup() *xdominion.XTable {
	t := xdominion.NewXTable("adminmenu_group", "adminmenu_grp_")

	// Key of the accessgroup
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	// Name of the group.
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
