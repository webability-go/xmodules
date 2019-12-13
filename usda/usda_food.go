package usda

import (
  "github.com/webability-go/xdominion"
)

/*
  Tabla de alimentos de la USDA
*/

func usda_food() *xdominion.XTable {
  t := xdominion.NewXTable("usda_food", "usda_fo_")

  t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 5, Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.PK},
                                                    } })   // PK

  t.AddField(xdominion.XFieldVarChar{Name: "group", Size: 4, Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                                xdominion.XConstraint{Type: xdominion.FK, Data: []string{"usda_group", "usda_gr_key"}},
                                                    } })

  t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 200, Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                    } })

  t.AddField(xdominion.XFieldVarChar{Name: "abbr", Size: 60, Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                    } })

  t.AddField(xdominion.XFieldVarChar{Name: "commonname", Size: 100 })

  t.AddField(xdominion.XFieldFloat{Name: "nfactor" })

  t.AddField(xdominion.XFieldFloat{Name: "profactor" })

  t.AddField(xdominion.XFieldFloat{Name: "fatfactor" })

  t.AddField(xdominion.XFieldFloat{Name: "chofactor" })

  t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                    } })

  return t
}

