package log

import (
	"github.com/webability-go/xdominion"
)

/*
  Main log table for data users
*/

func log_log() *xdominion.XTable {
	t := xdominion.NewXTable("log_log", "log_log_")
	// CLAVE
	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // AI, PK

	// CHEF
	t.AddField(xdominion.XFieldInteger{Name: "user", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.FK, Data: []string{"user_user", "user_usr_key"}},
		xdominion.XConstraint{Type: xdominion.NN},
	}}) // NN

	// OBJETO
	t.AddField(xdominion.XFieldVarChar{Name: "object", Size: 40, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.MI, Data: []string{"externkey"}},
	}}) //UI

	// ACCION
	t.AddField(xdominion.XFieldVarChar{Name: "action", Size: 4000})

	// TIMESTAMP
	t.AddField(xdominion.XFieldDateTime{Name: "timestamp"})

	// CLAVEEXTERNA
	t.AddField(xdominion.XFieldVarChar{Name: "externkey", Size: 40})

	return t

}
