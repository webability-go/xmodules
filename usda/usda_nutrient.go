package usda

import (
  "github.com/webability-go/xdominion"
)

/*
  Tabla de nutrientes de la USDA
*/

func usda_nutrient() *xdominion.XTable {
  t := xdominion.NewXTable("usda_nutrient", "usda_nu_")

  t.AddField(xdominion.XFieldVarChar{Name: "key", Size: 3, Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.PK},
                                                    } })   // PK

  t.AddField(xdominion.XFieldVarChar{Name: "unit", Size: 7, Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                    } })

  t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 60, Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                    } })

  t.AddField(xdominion.XFieldVarChar{Name: "tag", Size: 20 })

  t.AddField(xdominion.XFieldInteger{Name: "decimal", Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                    } })

  t.AddField(xdominion.XFieldInteger{Name: "order", Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                    } })

  t.AddField(xdominion.XFieldDateTime{Name: "lastmodif", Constraints: xdominion.XConstraints{
                                                                xdominion.XConstraint{Type: xdominion.NN},
                                                    } })

  return t
}

