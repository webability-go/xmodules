package ingredient

import (
  "github.com/webability-go/xdominion"
)

/*
  Clasificaciones de recetas
*/

func kl_ingrediente() *xdominion.XTable {
  t := xdominion.NewXTable("kl_ingrediente", "kl_ing_")

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

  // pasillo
  t.AddField(xdominion.XFieldInteger{Name: "pasilloingrediente", Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_ingredientepasillo", "kl_pas_clave"}},
                                                    xdominion.XConstraint{Type: xdominion.IN},
                                                    xdominion.XConstraint{Type: xdominion.NN},
                                                 } })

  // padre
  t.AddField(xdominion.XFieldInteger{Name: "padre", Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_ingrediente", "kl_ing_clave"}},
                                                 } })

  // usda
  t.AddField(xdominion.XFieldInteger{Name: "usda", Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.FK, Data: []string{"usda_food", "usda_fo_key"}},
                                                 } })

  // Unidad SI para pieza
  t.AddField(xdominion.XFieldInteger{Name: "unidadsi", Constraints: xdominion.XConstraints{
                                                    xdominion.XConstraint{Type: xdominion.FK, Data: []string{"kl_medida", "kl_med_clave"}},
                                                 } })

  // Cantidad en unidad SI para pieza
  t.AddField(xdominion.XFieldFloat{Name: "cantidad" })

  // densidad si volumen
  t.AddField(xdominion.XFieldFloat{Name: "densidad" })

  return t
}
