package metric

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

const (
	reFloat    = "^[0-9.]+$"
	reFraction = "^(([0-9]+)\\s+){0,1}([0-9]+)/([0-9]+)$"

	SI_GRAM  = 42
	SI_LITER = 49
	SI_METER = 205

	SI_KILO       = 45
	SI_CENTILITER = 14
	SI_MILILITER  = 55

	// solido
	SA_OUNCE = 56
	SA_POUND = 48

	// liquido
	SA_LOQUIDOUNCE = 207
	SA_GALLON      = 39

	// solido o liquido, aplicar densidad si necesario
	SC_CUP      = 86
	SC_SPOON    = 29
	SC_TEASPOON = 30

	UNIT_NOTVISIBLE = 102 // pieza sin escribir el nombre
	UNIT_VISIBLE    = 62  // pieza escribiendo el nombre
)

func GetMetric(sitecontext *context.Context, clave int, lang language.Tag) *StructureMetric {

	canonical := lang.String()

	data, _ := sitecontext.GetCache("metrics:" + canonical).Get(strconv.Itoa(clave))
	if data == nil {
		sm := CreateStructureMetricByKey(sitecontext, clave, lang)
		if sm == nil {
			sitecontext.Log("graph", "xmodules::metrics::GetMetric: No hay medida creada:", clave, lang)
			return nil
		}
		sitecontext.GetCache("metric:"+canonical).Set(strconv.Itoa(clave), sm)
		return sm.(*StructureMetric)
	}
	return data.(*StructureMetric)
}

// 1. si SOLO tiene 0-9 y . : entonces es int/float: convierte directo
// 2. si tiene SOLO con anteriores / y ' ' puede ser una fraccion
// si tiene otra cosa: no podemos convertir (letras, guiones, espacios etc)
// Regresa 0 si la cantidad es 0 o vacia
// regresa -1 si la cantidad no pudo ser comprendida (por ejemplo, puras letras, o un float con un formato erroneo)
func ParseQuantity(cantidad string) float64 {

	cantidad = strings.TrimSpace(cantidad)
	if cantidad == "" {
		return 0
	}
	fl := regexp.MustCompile(reFloat)
	fr := regexp.MustCompile(reFraction)

	if fl.MatchString(cantidad) {
		f, err := strconv.ParseFloat(cantidad, 64)
		if err == nil {
			return f
		}
	} else if fr.MatchString(cantidad) {
		// "1/4" o "1 1/4" mas de un espacio, error
		// solo una /, en el ultimo bloque, sino error
		// divide si ok el ultimo, sino error
		// suma los dos
		data := fr.FindAllStringSubmatch(cantidad, -1)
		var v float64 = 0
		var q float64 = 0
		var d float64 = 0
		if data[0][2] != "" {
			v, _ = strconv.ParseFloat(data[0][2], 64)
		}
		if data[0][3] != "" {
			q, _ = strconv.ParseFloat(data[0][3], 64)
		}
		if data[0][4] != "" {
			d, _ = strconv.ParseFloat(data[0][4], 64)
		}
		if q != 0 && d != 0 {
			v += q / d
		}
		return v
	}
	return -1
}

// system is 1 - "si", 2 - "sa" or 3 - "sc" (campo "tipo" de la tabla)
// 0 - as captured
// 1 - si = sistema internacional
// 2 - sa = sistema americano
// 3 - sc = sistema de cocina
// density is needed for weight-volume conversions (approx 0.7 to 2.5 generally)
// state is the type of metric we prefer for this ingredient: (basically for SI and SA), since SC is all volume
// 1 = solid
// 2 = liquid
func ConvertMetrics(sitecontext *context.Context, quantity float64, density float64, state int, fromMetric *xdominion.XRecord, system int) (float64, *xdominion.XRecord) {

	fromsystem, _ := fromMetric.GetInt("type")
	if fromsystem == system || fromsystem > 3 || system == 0 { // nothing to change, no calculable o el mismo sistema
		return quantity, fromMetric
	}

	fmt.Println("Vamos a calcular:", system, quantity, density, fromMetric)

	unidadsi, _ := fromMetric.GetInt("isunit")
	factor, _ := fromMetric.GetFloat("factorconversion")
	toMetric := GetMetric(sitecontext, unidadsi, language.Spanish) // only numbers interest us here

	// busca la medida más adhoc en el sistema correspondeiente
	if system == 1 { // vamos a usar el sistema de conversión directo de fromMetric
		// convierte a SI
		toquantity := quantity * factor
		if state == 1 && unidadsi == SI_LITER { // queremos solido pero es volumen, aplicar densidad
			toquantity = toquantity * density
			unidadsi = SI_GRAM
			toMetric = GetMetric(sitecontext, unidadsi, language.Spanish) // only numbers interest us here
		} else if state == 2 && unidadsi == SI_GRAM { // queremos volumen, es solido, aplicar 1/densidad
			toquantity = toquantity / density
			unidadsi = SI_LITER
			toMetric = GetMetric(sitecontext, unidadsi, language.Spanish) // only numbers interest us here
		} else {
			toMetric = GetMetric(sitecontext, unidadsi, language.Spanish) // only numbers interest us here
		}

		// buscamos el metrico más conveniente
		if unidadsi == SI_GRAM {
			tometric := GetMetric(sitecontext, SI_KILO, language.Spanish).GetData()
			factor, _ := tometric.GetFloat("factorconversion")
			newquantity := toquantity * factor
			if newquantity >= 1 {
				return newquantity, tometric
			}
		}
		// check si 1/4, 1/2, 3/4, 1/3, 2/3 de litro

		// TODO(phil) escanear las metricas de mismo tipo, convertir y tomar la cantidad más pequeña > 1

		return toquantity, toMetric.GetData()
	}

	// de sc o sa a sa o sc: convertimos primero a si
	var toquantity float64
	if fromsystem != 1 {
		// convierte a SI
		toquantity = quantity * factor
	} else {
		toquantity = quantity
	}

	// conversion final
	if system == 2 { // convert a SA: prueba los varios posibles
		if unidadsi == SI_GRAM {
			// intentar SA_ONZA, SA_LIBRA
			return quantity, fromMetric
		} else if unidadsi == SI_LITER {
			// intentar SA_GALON, SA_ONZALIQUIDA
			return quantity, fromMetric
		}
		// no se puede convertir !! guardamos el original
		return quantity, fromMetric
	}

	// system == 3
	if unidadsi == SI_GRAM {
		// intentar SC_TAZA, SC_CUCHARADA, SC_CUCHARADIDA con densidad

	} else if unidadsi == SI_LITER {
		// intentar SC_TAZA, SC_CUCHARADA, SC_CUCHARADIDA sin densidad
		tometric := GetMetric(sitecontext, SC_CUP, language.Spanish).GetData()
		factor, _ := tometric.GetFloat("factorconversion")
		newquantity := toquantity * factor
		if newquantity >= 1 {
			return newquantity, tometric
		}
		tometric = GetMetric(sitecontext, SC_SPOON, language.Spanish).GetData()
		factor, _ = tometric.GetFloat("factorconversion")
		newquantity = toquantity * factor
		if newquantity >= 1 {
			return newquantity, tometric
		}
		tometric = GetMetric(sitecontext, SC_TEASPOON, language.Spanish).GetData()
		factor, _ = tometric.GetFloat("factorconversion")
		newquantity = toquantity * factor
		return newquantity, tometric
	}

	// no se puede convertir !! guardamos el original
	return quantity, fromMetric
}
