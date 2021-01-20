package adminmenu

import (
	"github.com/webability-go/xdominion"
)

// TABLE: adminmenu_option:
// The menu tree for administration
func adminmenuOption() *xdominion.XTable {
	t := xdominion.NewXTable("adminmenu_option", "adminmenu_opt_")

	// Key of the accessgroup
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	// Group (hierarchy)
	t.AddField(xdominion.XFieldVarChar{Name: "group", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"adminmenu_group", "adminmenu_grp_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	t.AddField(xdominion.XFieldInteger{Name: "order", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	// Father (hierarchy)
	t.AddField(xdominion.XFieldVarChar{Name: "father", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"adminmenu_option", "adminmenu_opt_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	// Name of the group.
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	t.AddField(xdominion.XFieldVarChar{Name: "access", Size: 30, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_access", "user_acc_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	// what to call on click on name & icon 1
	t.AddField(xdominion.XFieldVarChar{Name: "icon1", Size: 255})
	t.AddField(xdominion.XFieldVarChar{Name: "call1", Size: 255})
	t.AddField(xdominion.XFieldVarChar{Name: "description1", Size: 4000})

	// what to call on click on icon 2 (first icon after)
	t.AddField(xdominion.XFieldVarChar{Name: "icon2", Size: 255})
	t.AddField(xdominion.XFieldVarChar{Name: "call2", Size: 255})
	t.AddField(xdominion.XFieldVarChar{Name: "description2", Size: 4000})

	// what to call on click on icon 3 (second icon after)
	t.AddField(xdominion.XFieldVarChar{Name: "icon3", Size: 255})
	t.AddField(xdominion.XFieldVarChar{Name: "call3", Size: 255})
	t.AddField(xdominion.XFieldVarChar{Name: "description3", Size: 4000})

	return t
}
