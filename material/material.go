package material

import (
  "strconv"

  "golang.org/x/text/language"

  "xmodules/context"
  "xmodules/metrics"
)

func InitMaterial(sitecontext *context.Context, databasename string) error {

  buildTables(sitecontext, databasename)
  buildCache(sitecontext)

  return nil
}

func GetMaterial(sitecontext *context.Context, clave int, lang language.Tag) *StructureMaterial {

	canonical := lang.String()

  data, _ := sitecontext.Caches["materiales:" + canonical].Get(strconv.Itoa(clave))
  if data == nil {
    sm := CreateStructureMaterialByKey(sitecontext, clave, lang)
    if sm == nil {
      sitecontext.Logs["graph"].Println("xmodules::material::GetMaterial: No hay material creado:", clave, lang)
      return nil
    }
    sitecontext.Caches["materiales:" + canonical].Set(strconv.Itoa(clave), sm)
		return sm.(*StructureMaterial)
  }
  return data.(*StructureMaterial)
}


func GetMaterialCompositeName(sitecontext *context.Context, quantity string, materialkey int, metrickey int, extra string, system int, lang language.Tag) string {

  materialstructure := GetMaterial(sitecontext, materialkey, lang)
  metricstructure := metrics.GetMetric(sitecontext, metrickey, lang)
  if materialstructure == nil || metricstructure == nil {
    return extra
  }
  materialdata := materialstructure.GetData()
  metricdata := metricstructure.GetData()
  if materialdata == nil || metricdata == nil {
    return extra
  }

  xquantity := metrics.ParseQuantity(quantity)
  matnamesingular, _ := materialdata.GetString("nombre")
  matnameplural, _ := materialdata.GetString("plural")
  metricingular, _ := metricdata.GetString("nombre")
  metricplural, _ := metricdata.GetString("plural")

  composite := quantity + " "

  // pieza sin decirlo = unidades del ingrediente
  if metrickey == metrics.UNIT_NOTVISIBLE {
    if xquantity == 1.0 {
      composite += matnamesingular
    } else {
      composite += matnameplural
    }
  } else {
    if xquantity == 1.0 {
      composite += metricingular
    } else {
      composite += metricplural
    }

    composite += " "
    composite += matnamesingular
  }

  if extra != "" {
    composite += ", " + extra
  }

  return composite
}
