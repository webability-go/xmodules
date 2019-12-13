package material

import (
  "github.com/webability-go/xdominion"
)

/*
  Clasificaciones de recetas
*/

func kl_material() *xdominion.XTable {
  t := xdominion.NewXTable("kl_material", "kl_mat_")

  t.AddField(xdominion.XFieldInteger{Name: "clave", Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.PK},
                                                    xdominion.XConstraint{Type: xdominion.AI},
                                                 } })   // AI, PK

  // nombre
  t.AddField(xdominion.XFieldVarChar{Name: "nombre", Size: 255, Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.NN},
                                                 } })

  // plural
  t.AddField(xdominion.XFieldVarChar{Name: "plural", Size: 255, Constraints: xdominion.XConstraints{
                                                xdominion.XConstraint{Type: xdominion.NN},
                                              } })

  // fecha last modif
  t.AddField(xdominion.XFieldDateTime{Name: "fecha", Constraints: xdominion.XConstraints{
                                                xdominion.XConstraint{Type: xdominion.NN},
                                              } })

  return t
}
