package material

import (
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/metric"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID         = "material"
	VERSION          = "1.0.0"
	TRANSLATIONTHEME = "material"
)

func InitModule(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	go buildCache(sitecontext)
	sitecontext.Modules[MODULEID] = VERSION

	return nil
}

func SynchronizeModule(sitecontext *context.Context) []string {

	translation.AddTheme(sitecontext, TRANSLATIONTHEME, "Material", translation.SOURCETABLE, "", "name,plural")

	messages := []string{}
	messages = append(messages, "Analysing material_material table.")
	num, err := sitecontext.Tables["material_material"].Count(nil)
	if err != nil || num == 0 {
		err1 := sitecontext.Tables["material_material"].Synchronize()
		if err1 != nil {
			messages = append(messages, "The table material_material was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table material_material was created (again)")
		}
	} else {
		messages = append(messages, "The table material_material was not created because it contains data.")
	}

	// fill metric and translations
	return messages
}

func GetMaterial(sitecontext *context.Context, clave int, lang language.Tag) *StructureMaterial {

	canonical := lang.String()

	data, _ := sitecontext.Caches["materiales:"+canonical].Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructureMaterialByKey(sitecontext, clave, lang)
		if sm == nil {
			sitecontext.Logs["graph"].Println("xmodules::material::GetMaterial: No hay material creado:", clave, lang)
			return nil
		}
		sitecontext.Caches["materiales:"+canonical].Set(strconv.Itoa(clave), sm)
		return sm.(*StructureMaterial)
	}
	return data.(*StructureMaterial)
}

func GetMaterialCompositeName(sitecontext *context.Context, quantity string, materialkey int, metrickey int, extra string, system int, lang language.Tag) string {

	materialstructure := GetMaterial(sitecontext, materialkey, lang)
	metricstructure := metric.GetMetric(sitecontext, metrickey, lang)
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
