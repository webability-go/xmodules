package ingredient

import (
  "github.com/webability-go/xdominion"
)

/*
  Clasificaciones de recetas
*/

func kl_ingredientepasillo() *xdominion.XTable {
  t := xdominion.NewXTable("kl_ingredientepasillo", "kl_pas_")

  t.AddField(xdominion.XFieldInteger{Name: "clave", Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.PK},
                                                    xdominion.XConstraint{Type: xdominion.AI},
                                                 } })   // AI, PK

  // nombre
  t.AddField(xdominion.XFieldVarChar{Name: "nombre", Size: 50, Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.NN},
                                                 } })

  // orden
  t.AddField(xdominion.XFieldInteger{Name: "orden", Constraints: xdominion.XConstraints{
                                                   xdominion.XConstraint{Type: xdominion.NN},
                                                } })

  // fecha last modif
  t.AddField(xdominion.XFieldDateTime{Name: "fecha", Constraints: xdominion.XConstraints{
                                                  xdominion.XConstraint{Type: xdominion.NN},
                                                } })

  return t
}
