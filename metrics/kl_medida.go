package metrics

import (
  "github.com/webability-go/xdominion"
)

/*
  Medidas de ingredientes
*/

func kl_medida() *xdominion.XTable {
  t := xdominion.NewXTable("kl_medida", "kl_med_")

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

  // 1 = SI, 2 = US, 3 = cócina, 4 = insignificante, 5 = no cuantificable
  t.AddField(xdominion.XFieldInteger{Name: "tipo"})

  // 1 si es medida SI oficial (metro, litro, gramo, galón, libra, taza... )
  t.AddField(xdominion.XFieldInteger{Name: "oficial"})

  // unidad sistema internacional
  t.AddField(xdominion.XFieldInteger{Name: "unidadsi", Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_medida", "kl_med_clave"}},
                                                 } })

  // unidad sistema internacional
  t.AddField(xdominion.XFieldFloat{Name: "factorconversion" })

  return t
}

