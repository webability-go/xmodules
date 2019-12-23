package user

import (
	"github.com/webability-go/xdominion"
)

// TABLE: user_session:
// All the sessions of the usesr
func userSession() *xdominion.XTable {
	t := xdominion.NewXTable("user_session", "user_ses_")

	// Key (CID) of the connection
	t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 36, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // PK

	// user
	t.AddField(xdominion.XFieldInteger{Name: "user", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_user", "user_usr_key"}},
		xdominion.XConstraint{Type: xdominion.IN},
		xdominion.XConstraint{Type: xdominion.DC},
	}})

	// date of creation
	t.AddField(xdominion.XFieldDateTime{Name: "creationdate", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// date of last modification
	t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// IP of the connection IPV4, IPV6, Chain of IPs
	t.AddField(xdominion.XFieldVarChar{Name: "ip", Size: 40, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// 1 if long login (6 months validity)
	t.AddField(xdominion.XFieldInteger{Name: "longlogin"})

	// origin
	t.AddField(xdominion.XFieldVarChar{Name: "origin", Size: 20, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// device
	t.AddField(xdominion.XFieldVarChar{Name: "device", Size: 10, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	return t
}
