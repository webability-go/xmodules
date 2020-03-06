// Modules dependency:
// Necesita USDA y METRICS para funcionar
package ingredient

import (
	"fmt"
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/context"
	"github.com/webability-go/xmodules/metric"
	"github.com/webability-go/xmodules/translation"
)

const (
	MODULEID              = "ingredient"
	VERSION               = "2.0.0"
	TRANSLATIONTHEME      = "ingredient"
	TRANSLATIONTHEMEAISLE = "ingaisle"
)

func InitModule(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	createCache(sitecontext)
	sitecontext.SetModule(MODULEID, VERSION)

	go buildCache(sitecontext)

	return nil
}

func SynchronizeModule(sitecontext *context.Context) []string {

	translation.AddTheme(sitecontext, TRANSLATIONTHEME, "Ingredients", translation.SOURCETABLE, "", "name,plural")
	translation.AddTheme(sitecontext, TRANSLATIONTHEMEAISLE, "Ingredients Aisles", translation.SOURCETABLE, "", "name")

	messages := []string{}

	messages = append(messages, "Analysing ingredient_aisle table.")
	num, err := sitecontext.GetTable("ingredient_aisle").Count(nil)
	if err != nil || num == 0 {
		err1 := sitecontext.GetTable("ingredient_aisle").Synchronize()
		if err1 != nil {
			messages = append(messages, "The table ingredient_aisle was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table ingredient_aisle was created (again)")
		}
	} else {
		messages = append(messages, "The table ingredient_aisle was not created because it contains data.")
	}

	messages = append(messages, "Analysing ingredient_ingredient table.")
	num, err = sitecontext.GetTable("ingredient_ingredient").Count(nil)
	if err != nil || num == 0 {
		err1 := sitecontext.GetTable("ingredient_ingredient").Synchronize()
		if err1 != nil {
			messages = append(messages, "The table ingredient_ingredient was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table ingredient_ingredient was created (again)")
		}
	} else {
		messages = append(messages, "The table ingredient_ingredient was not created because it contains data.")
	}

	// fill metric and translations
	return messages
}

func GetPasillo(sitecontext *context.Context, clave int, lang language.Tag) *StructurePasillo {

	canonical := lang.String()

	data, _ := sitecontext.GetCache("ingredient:pasillos:" + canonical).Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructurePasilloByKey(sitecontext, clave, lang)
		if sm == nil {
			sitecontext.Log("graph", "xmodules::ingredient::GetPasillo: No hay pasillo creado:", clave, lang)
			return nil
		}
		sitecontext.GetCache("ingredient:pasillos:"+canonical).Set(strconv.Itoa(clave), sm)
		return sm.(*StructurePasillo)
	}
	return data.(*StructurePasillo)
}

func GetIngredient(sitecontext *context.Context, clave int, lang language.Tag) *StructureIngredient {

	canonical := lang.String()

	data, _ := sitecontext.GetCache("ingredient:ingredientes:" + canonical).Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructureIngredientByKey(sitecontext, clave, lang)
		if sm == nil {
			sitecontext.Log("graph", "xmodules::ingredient::GetIngredient: No hay ingrediente creado:", clave, lang)
			return nil
		}
		sitecontext.GetCache("ingredient:ingredientes:"+canonical).Set(strconv.Itoa(clave), sm)
		return sm.(*StructureIngredient)
	}
	return data.(*StructureIngredient)
}

func ConvertToSI(sitecontext *context.Context, ingrediente int, scantidad string, medida int, cantidadsi int, medidasi int) (float64, int, string) {

	if scantidad == "" || medida == 0 {
		return 0, -1, "Los parámetros aún son incompletos"
	}

	if cantidadsi != 0 && medidasi != 0 {
		return ConvertToSI(sitecontext, ingrediente, fmt.Sprint(cantidadsi), medidasi, 0, 0)
	}

	// cantidad es correcta ?
	cantidad := metric.ParseQuantity(scantidad)
	if cantidad < 0 {
		return 0, -2, "La cantidad no se pudo interpretar"
	}

	var factorconversion float64 = 0
	// medida es cuantificable ?
	medidadata := metric.GetMetric(sitecontext, medida, language.Spanish) // el idioma no importa aqui, estamos calculando cifras solamente
	unidadsi, _ := medidadata.Data.GetInt("unidadsi")
	if unidadsi != 0 {
		factorconversion, _ = medidadata.Data.GetFloat("factorconversion")
	} else {
		// CASO ESPECIAL: por pieza
		// verificar si tenemos un peso/unidad en el ingrediente x defecto (ejemplo: 1 huevo, 1 manzana, etc)
		if medida == metric.UNIT_NOTVISIBLE || medida == metric.UNIT_VISIBLE {
			ingredientedata := GetIngredient(sitecontext, ingrediente, language.Spanish) // el idioma no importa aqui, estamos calculando cifras solamente
			usi, _ := ingredientedata.Data.GetInt("unidadsi")
			if usi != 0 {
				unidadsi = usi
				csi, _ := ingredientedata.Data.GetFloat("cantidad")
				medidadata = metric.GetMetric(sitecontext, usi, language.Spanish)
				factorconversion, _ = medidadata.Data.GetFloat("factorconversion")
				cantidad = cantidad * csi
			} else {
				return 0, -4, "Falta una unidad de conversión en Sistema Internacional de la pieza"
			}
		} else {
			return 0, -3, "Falta una unidad y cantidad de conversión en Sistema Internacional"
		}
	}

	return factorconversion * cantidad, unidadsi, ""
}

func GetIngredientCompositeName(sitecontext *context.Context, quantity string, ingredientkey int, metrickey int, extra string, system int, lang language.Tag) string {

	nombrecompuesto := ""
	switch lang {
	case language.Spanish:
		nombrecompuesto = CompositeNameSpanish(sitecontext, quantity, ingredientkey, metrickey, extra, system)
	case language.English:
		nombrecompuesto = CompositeNameEnglish(sitecontext, quantity, ingredientkey, metrickey, extra, system)
	}
	return nombrecompuesto
}

func CompositeNameSpanish(sitecontext *context.Context, quantity string, ingredientkey int, metrickey int, extra string, system int) string {

	xquantity := metric.ParseQuantity(quantity)
	ingredient := GetIngredient(sitecontext, ingredientkey, language.Spanish)
	if ingredient == nil {
		// log this
		return "Error in CompositeNameSpanish::ingredient NIL"
	}
	ingredientdata := ingredient.GetData()
	metricstructure := metric.GetMetric(sitecontext, metrickey, language.Spanish)
	if metricstructure == nil {
		// log this
		return "Error in CompositeNameSpanish::metric NIL"
	}
	metricdata := metricstructure.GetData()

	density, _ := ingredientdata.GetFloat("density")
	state, _ := ingredientdata.GetInt("type")

	if system != 0 { // 1, 2, 3: convertir al sistema solicitado antes de crear composite
		xquantity, metricdata = metric.ConvertMetrics(sitecontext, xquantity, density, state, metricdata, system)
		quantity = fmt.Sprint(xquantity)
	}

	ingnamesingular, _ := ingredientdata.GetString("nombre")
	ingnameplural, _ := ingredientdata.GetString("plural")
	metricsingular, _ := metricdata.GetString("nombre")
	metricplural, _ := metricdata.GetString("plural")

	composite := quantity + " "

	// pieza sin decirlo = unidades del ingrediente
	if metrickey == metric.UNIT_NOTVISIBLE {
		if xquantity == 1.0 {
			composite += ingnamesingular
		} else {
			composite += ingnameplural
		}
	} else {
		if xquantity == 1.0 {
			composite += metricsingular
		} else {
			composite += metricplural
		}

		composite += " de "

		composite += ingnamesingular
		// Podria ser plural tambien, por ejemplo 1 caja de galletas: VER REGLA
		//    compuesto += ingnameplural
	}

	if extra != "" {
		composite += ", " + extra
	}

	return composite
}

func CompositeNameEnglish(sitecontext *context.Context, quantity string, ingredientkey int, metrickey int, extra string, system int) string {

	xquantity := metric.ParseQuantity(quantity)
	ingredient := GetIngredient(sitecontext, ingredientkey, language.English)
	if ingredient == nil {
		// log this
		return "Error in CompositeNameSpanish::ingredient NIL"
	}
	ingredientdata := ingredient.GetData()
	metricstructure := metric.GetMetric(sitecontext, metrickey, language.English)
	if metricstructure == nil {
		// log this
		return "Error in CompositeNameSpanish::metric NIL"
	}
	metricdata := metricstructure.GetData()

	density, _ := ingredientdata.GetFloat("density")
	state, _ := ingredientdata.GetInt("type")

	if system != 0 { // 1, 2, 3: convertir al sistema solicitado antes de crear composite
		xquantity, metricdata = metric.ConvertMetrics(sitecontext, xquantity, density, state, metricdata, system)
		quantity = fmt.Sprint(xquantity)
	}

	ingnamesingular, _ := ingredientdata.GetString("name")
	ingnameplural, _ := ingredientdata.GetString("plural")
	metricsingular, _ := metricdata.GetString("name")
	metricplural, _ := metricdata.GetString("plural")

	composite := quantity + " "

	// pieza sin decirlo = unidades del ingrediente
	if metrickey == metric.UNIT_NOTVISIBLE {
		if xquantity == 1.0 {
			composite += ingnamesingular
		} else {
			composite += ingnameplural
		}
	} else {
		if xquantity == 1.0 {
			composite += metricsingular
		} else {
			composite += metricplural
		}

		composite += " "

		composite += ingnamesingular
		// Podria ser plural tambien, por ejemplo 1 caja de galletas: VER REGLA
		//    compuesto += ingnameplural
	}

	if extra != "" {
		composite += ", " + extra
	}

	return composite
}
