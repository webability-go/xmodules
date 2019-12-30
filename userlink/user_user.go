package userlink

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_user:
// All the administrators of the base system, on a remote system. Those are the data copied from the controller system
func userUser() *xdominion.XTable {
	t := xdominion.NewXTable("user_user", "user_usr_")

	// Key of the user, automatic consecutive
	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	// The status of the administrator:
	// X = Suspended, A = Active, S = Superuser  (superusers can connect to master and have all the rights by default)
	t.AddField(xdominion.XFieldVarChar{Name: "status", Size: 1, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// Name of the administrator.
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 200, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// Mail of the administrator
	t.AddField(xdominion.XFieldVarChar{Name: "mail", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
