package material

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xmodules/metric"
)

func GetMaterial(ds applications.Datasource, clave int, lang language.Tag) *StructureMaterial {

	canonical := lang.String()

	data, _ := ds.GetCache("materiales:" + canonical).Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructureMaterialByKey(ds, clave, lang)
		if sm == nil {
			ds.Log("graph", "xmodules::material::GetMaterial: No hay material creado:", clave, lang)
			return nil
		}
		ds.GetCache("materiales:"+canonical).Set(strconv.Itoa(clave), sm)
		return sm.(*StructureMaterial)
	}
	return data.(*StructureMaterial)
}

func GetMaterialCompositeName(ds applications.Datasource, quantity string, materialkey int, metrickey int, extra string, system int, lang language.Tag) string {

	materialstructure := GetMaterial(ds, materialkey, lang)
	metricstructure := metric.GetMetric(ds, metrickey, lang)
	if materialstructure == nil || metricstructure == nil {
		return extra
	}
	materialdata := materialstructure.GetData()
	metricdata := metricstructure.GetData()
	if materialdata == nil || metricdata == nil {
		return extra
	}

	xquantity := metric.ParseQuantity(quantity)
	matnamesingular, _ := materialdata.GetString("nombre")
	matnameplural, _ := materialdata.GetString("plural")
	metricingular, _ := metricdata.GetString("nombre")
	metricplural, _ := metricdata.GetString("plural")

	composite := quantity + " "

	// pieza sin decirlo = unidades del ingrediente
	if metrickey == metric.UNIT_NOTVISIBLE {
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
