package stat

import (
	"github.com/webability-go/xdominion"
)

/*
  Client parameters table
*/

func stat_stat(prefix string, num string) *xdominion.XTable {
	t := xdominion.NewXTable(prefix+"stat_"+num, prefix+"sta_")

	t.AddField(xdominion.XFieldInteger{Name: "key", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // AI, PK

	// skin
	t.AddField(xdominion.XFieldVarChar{Name: "skin", Size: 30})

	// pagina
	t.AddField(xdominion.XFieldVarChar{Name: "page", Size: 250, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	// uri
	t.AddField(xdominion.XFieldVarChar{Name: "uri", Size: 250})

	// ip
	t.AddField(xdominion.XFieldVarChar{Name: "ip", Size: 40, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// host
	t.AddField(xdominion.XFieldVarChar{Name: "host", Size: 255})

	// fecha
	t.AddField(xdominion.XFieldDateTime{Name: "date", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
		xdominion.XConstraint{Type: xdominion.IN},
	}})

	// Tiempo
	t.AddField(xdominion.XFieldFloat{Name: "time", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// params
	t.AddField(xdominion.XFieldText{Name: "params", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})

	// method
	t.AddField(xdominion.XFieldVarChar{Name: "method", Size: 6})

	// code
	t.AddField(xdominion.XFieldInteger{Name: "code"})

	// server
	t.AddField(xdominion.XFieldVarChar{Name: "server", Size: 15})

	// client
	t.AddField(xdominion.XFieldInteger{Name: "client"})

	// dispositivo
	t.AddField(xdominion.XFieldVarChar{Name: "device", Size: 15})

	// browser id
	t.AddField(xdominion.XFieldVarChar{Name: "browser", Size: 4000})

	return t
}
