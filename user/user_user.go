package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_user:
// All the administrators of the base system.
func userUser() *xdominion.XTable {
	t := xdominion.NewXTable("user_usr", "user_usr_")

	// Key of the user, automatic consecutive
	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
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

	// Username (login) of the administrator
	t.AddField(xdominion.XFieldVarChar{Name: "un", Size: 200, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.UI},
	}})

	// Password of the administrator
	// natural if "not encrypted", or MD5 if the flag is "encrypted"
	// If the flag is encrypted, recuperating the password is assigning any random password
	t.AddField(xdominion.XFieldVarChar{Name: "pw", Size: 200, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// Mail of the administrator
	t.AddField(xdominion.XFieldVarChar{Name: "mail", Size: 255, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// sex, M/F
	t.AddField(xdominion.XFieldVarChar{Name: "sex", Size: 1, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// fields of special data of user
	// example: BOSS=...\nTEL=...\nUNIT=...
	// Those fields are used into the data of the user
	t.AddField(xdominion.XFieldText{Name: "fields"})

	// extra info
	t.AddField(xdominion.XFieldText{Name: "info"})

	// The administrator who create or is responsible of this one.
	// if null, it is a super user
	t.AddField(xdominion.XFieldInteger{Name: "father", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_user", "user_usr_key"}},
	}})

	// date of creation
	t.AddField(xdominion.XFieldDateTime{Name: "creationdate", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// date of last modification
	t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
