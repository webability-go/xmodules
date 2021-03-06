// Modules dependency:
// Necesita USDA y METRICS para funcionar
package ingredient

import (
	"fmt"
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/metric"
)

func GetPasillo(ctx *base.Datasource, clave int, lang language.Tag) *StructurePasillo {

	canonical := lang.String()

	data, _ := ctx.GetCache("ingredient:pasillos:" + canonical).Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructurePasilloByKey(ctx, clave, lang)
		if sm == nil {
			ctx.Log("graph", "xmodules::ingredient::GetPasillo: No hay pasillo creado:", clave, lang)
			return nil
		}
		ctx.GetCache("ingredient:pasillos:"+canonical).Set(strconv.Itoa(clave), sm)
		return sm.(*StructurePasillo)
	}
	return data.(*StructurePasillo)
}

func GetIngredient(ctx *base.Datasource, clave int, lang language.Tag) *StructureIngredient {

	canonical := lang.String()

	data, _ := ctx.GetCache("ingredient:ingredientes:" + canonical).Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructureIngredientByKey(ctx, clave, lang)
		if sm == nil {
			ctx.Log("graph", "xmodules::ingredient::GetIngredient: No hay ingrediente creado:", clave, lang)
			return nil
		}
		ctx.GetCache("ingredient:ingredientes:"+canonical).Set(strconv.Itoa(clave), sm)
		return sm.(*StructureIngredient)
	}
	return data.(*StructureIngredient)
}

func ConvertToSI(ctx *base.Datasource, ingrediente int, scantidad string, medida int, cantidadsi int, medidasi int) (float64, int, string) {

	if scantidad == "" || medida == 0 {
		return 0, -1, "Los parámetros aún son incompletos"
	}

	if cantidadsi != 0 && medidasi != 0 {
		return ConvertToSI(ctx, ingrediente, fmt.Sprint(cantidadsi), medidasi, 0, 0)
	}

	// cantidad es correcta ?
	cantidad := metric.ParseQuantity(scantidad)
	if cantidad < 0 {
		return 0, -2, "La cantidad no se pudo interpretar"
	}

	var factorconversion float64 = 0
	// medida es cuantificable ?
	medidadata := metric.GetMetric(ctx, medida, language.Spanish) // el idioma no importa aqui, estamos calculando cifras solamente
	unidadsi, _ := medidadata.Data.GetInt("unidadsi")
	if unidadsi != 0 {
		factorconversion, _ = medidadata.Data.GetFloat("factorconversion")
	} else {
		// CASO ESPECIAL: por pieza
		// verificar si tenemos un peso/unidad en el ingrediente x defecto (ejemplo: 1 huevo, 1 manzana, etc)
		if medida == metric.UNIT_NOTVISIBLE || medida == metric.UNIT_VISIBLE {
			ingredientedata := GetIngredient(ctx, ingrediente, language.Spanish) // el idioma no importa aqui, estamos calculando cifras solamente
			usi, _ := ingredientedata.Data.GetInt("unidadsi")
			if usi != 0 {
				unidadsi = usi
				csi, _ := ingredientedata.Data.GetFloat("cantidad")
				medidadata = metric.GetMetric(ctx, usi, language.Spanish)
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

func GetIngredientCompositeName(ctx *base.Datasource, quantity string, ingredientkey int, metrickey int, extra string, system int, lang language.Tag) string {

	nombrecompuesto := ""
	switch lang {
	case language.Spanish:
		nombrecompuesto = CompositeNameSpanish(ctx, quantity, ingredientkey, metrickey, extra, system)
	case language.English:
		nombrecompuesto = CompositeNameEnglish(ctx, quantity, ingredientkey, metrickey, extra, system)
	}
	return nombrecompuesto
}

func CompositeNameSpanish(ctx *base.Datasource, quantity string, ingredientkey int, metrickey int, extra string, system int) string {

	xquantity := metric.ParseQuantity(quantity)
	ingredient := GetIngredient(ctx, ingredientkey, language.Spanish)
	if ingredient == nil {
		// log this
		return "Error in CompositeNameSpanish::ingredient NIL"
	}
	ingredientdata := ingredient.GetData()
	metricstructure := metric.GetMetric(ctx, metrickey, language.Spanish)
	if metricstructure == nil {
		// log this
		return "Error in CompositeNameSpanish::metric NIL"
	}
	metricdata := metricstructure.GetData()

	density, _ := ingredientdata.GetFloat("density")
	state, _ := ingredientdata.GetInt("type")

	if system != 0 { // 1, 2, 3: convertir al sistema solicitado antes de crear composite
		xquantity, metricdata = metric.ConvertMetrics(ctx, xquantity, density, state, metricdata, system)
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

func CompositeNameEnglish(ctx *base.Datasource, quantity string, ingredientkey int, metrickey int, extra string, system int) string {

	xquantity := metric.ParseQuantity(quantity)
	ingredient := GetIngredient(ctx, ingredientkey, language.English)
	if ingredient == nil {
		// log this
		return "Error in CompositeNameSpanish::ingredient NIL"
	}
	ingredientdata := ingredient.GetData()
	metricstructure := metric.GetMetric(ctx, metrickey, language.English)
	if metricstructure == nil {
		// log this
		return "Error in CompositeNameSpanish::metric NIL"
	}
	metricdata := metricstructure.GetData()

	density, _ := ingredientdata.GetFloat("density")
	state, _ := ingredientdata.GetInt("type")

	if system != 0 { // 1, 2, 3: convertir al sistema solicitado antes de crear composite
		xquantity, metricdata = metric.ConvertMetrics(ctx, xquantity, density, state, metricdata, system)
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
