// Modules dependency:
// Necesita USDA y METRICS para funcionar
package ingredient

import (
	"fmt"
	"strconv"

	"golang.org/x/text/language"

	"xmodules/context"
	"xmodules/metrics"
)

func InitIngredient(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	buildCache(sitecontext)

	return nil
}

func GetPasillo(sitecontext *context.Context, clave int, lang language.Tag) *StructurePasillo {

	canonical := lang.String()

	data, _ := sitecontext.Caches["ingredient:pasillos:"+canonical].Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructurePasilloByKey(sitecontext, clave, lang)
		if sm == nil {
			sitecontext.Logs["graph"].Println("xmodules::ingredient::GetPasillo: No hay pasillo creado:", clave, lang)
			return nil
		}
		sitecontext.Caches["ingredient:pasillos:"+canonical].Set(strconv.Itoa(clave), sm)
		return sm.(*StructurePasillo)
	}
	return data.(*StructurePasillo)
}

func GetIngredient(sitecontext *context.Context, clave int, lang language.Tag) *StructureIngredient {

	canonical := lang.String()

	data, _ := sitecontext.Caches["ingredient:ingredientes:"+canonical].Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructureIngredientByKey(sitecontext, clave, lang)
		if sm == nil {
			sitecontext.Logs["graph"].Println("xmodules::ingredient::GetIngredient: No hay ingrediente creado:", clave, lang)
			return nil
		}
		sitecontext.Caches["ingredient:ingredientes:"+canonical].Set(strconv.Itoa(clave), sm)
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
	cantidad := metrics.ParseQuantity(scantidad)
	if cantidad < 0 {
		return 0, -2, "La cantidad no se pudo interpretar"
	}

	var factorconversion float64 = 0
	// medida es cuantificable ?
	medidadata := metrics.GetMetric(sitecontext, medida, language.Spanish) // el idioma no importa aqui, estamos calculando cifras solamente
	unidadsi, _ := medidadata.Data.GetInt("unidadsi")
	if unidadsi != 0 {
		factorconversion, _ = medidadata.Data.GetFloat("factorconversion")
	} else {
		// CASO ESPECIAL: por pieza
		// verificar si tenemos un peso/unidad en el ingrediente x defecto (ejemplo: 1 huevo, 1 manzana, etc)
		if medida == metrics.UNIT_NOTVISIBLE || medida == metrics.UNIT_VISIBLE {
			ingredientedata := GetIngredient(sitecontext, ingrediente, language.Spanish) // el idioma no importa aqui, estamos calculando cifras solamente
			usi, _ := ingredientedata.Data.GetInt("unidadsi")
			if usi != 0 {
				unidadsi = usi
				csi, _ := ingredientedata.Data.GetFloat("cantidad")
				medidadata = metrics.GetMetric(sitecontext, usi, language.Spanish)
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

	xquantity := metrics.ParseQuantity(quantity)
	ingredientdata := GetIngredient(sitecontext, ingredientkey, language.Spanish).GetData()
	density, _ := ingredientdata.GetFloat("densidad")
	state, _ := ingredientdata.GetInt("tipo")
	metricdata := metrics.GetMetric(sitecontext, metrickey, language.Spanish).GetData()

	if system != 0 { // 1, 2, 3: convertir al sistema solicitado antes de crear composite
		xquantity, metricdata = metrics.ConvertMetrics(sitecontext, xquantity, density, state, metricdata, system)
		quantity = fmt.Sprint(xquantity)

		fmt.Println("Hemos calculado:", system, xquantity, metricdata)

	}

	ingnamesingular, _ := ingredientdata.GetString("nombre")
	ingnameplural, _ := ingredientdata.GetString("plural")
	metricsingular, _ := metricdata.GetString("nombre")
	metricplural, _ := metricdata.GetString("plural")

	composite := quantity + " "

	// pieza sin decirlo = unidades del ingrediente
	if metrickey == metrics.UNIT_NOTVISIBLE {
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

	xquantity := metrics.ParseQuantity(quantity)
	ingredientdata := GetIngredient(sitecontext, ingredientkey, language.English).GetData()
	metricdata := metrics.GetMetric(sitecontext, metrickey, language.English).GetData()
	ingnamesingular, _ := ingredientdata.GetString("nombre")
	ingnameplural, _ := ingredientdata.GetString("plural")
	metricsingular, _ := metricdata.GetString("nombre")
	metricplural, _ := metricdata.GetString("plural")

	composite := quantity + " "

	// pieza sin decirlo = unidades del ingrediente
	if metrickey == metrics.UNIT_NOTVISIBLE {
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
