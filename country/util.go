package country

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/translation"
)

// Order to load/synchronize tables:
var moduletablesorder = []string{
	"country_country",
}

// map[string] does not respect order
var moduletables = map[string]func() *xdominion.XTable{
	"country_country": countryCountry,
}

func buildTables(ctx *base.Datasource) {

	for _, tbl := range moduletablesorder {
		table := moduletables[tbl]()
		table.SetBase(ctx.GetDatabase())
		table.SetLanguage(language.English)
		ctx.SetTable(tbl, table)
	}
}

func createCache(ctx *base.Datasource) []string {

	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()
		ctx.SetCache("country:countries:"+canonical, xcore.NewXCache("country:countries:"+canonical, 0, 0))
	}
	return []string{}
}

func buildCache(ctx *base.Datasource) []string {

	// Lets protect us for race condition since map[] of Tables and XCaches are not thread safe
	country_country := ctx.GetTable("country_country")
	caches := map[string]*xcore.XCache{}
	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()
		caches["country:countries:"+canonical] = ctx.GetCache("country:countries:" + canonical)
	}

	// Loads all data in XCache
	countries, _ := country_country.SelectAll()
	if countries == nil {
		return []string{"No hay paises en la tabla"}
	}

	for _, lang := range ctx.GetLanguages() {
		canonical := lang.String()

		all := []string{}
		for _, m := range *countries {
			// creates structure on language
			str := CreateStructureCountryByData(ctx, m.Clone(), lang)
			key, _ := m.GetString("key")
			all = append(all, key)
			caches["country:countries:"+canonical].Set(key, str)
		}
		caches["country:countries:"+canonical].Set("all", all)
	}

	return []string{}
}

func loadTables(ctx *base.Datasource, filespath string) []string {

	// borra toda la data porque la vamos a insertar de nuevo (si se puede: FK bloquea)
	ctx.GetTable("country_country").Delete(nil)

	// 4 archivos de importación
	DMPCOUNTRY := filespath + "countries.en.dmp"

	num := 0
	data := readFile(DMPCOUNTRY)
	for _, r := range *data {
		key, _ := r.(*xdominion.XRecord).GetString("key")
		ctx.GetTable("country_country").Upsert(key, *r.(*xdominion.XRecord))
		num++
	}

	// insert into translation each country in spanish
	DMPCOUNTRY = filespath + "countries.es.dmp"
	data = readFile(DMPCOUNTRY)
	for _, r := range *data {
		key, _ := r.(*xdominion.XRecord).GetString("key")
		name, _ := r.(*xdominion.XRecord).GetString("name")

		err := translation.SetTranslation(ctx, name, TRANSLATIONTHEME, key, "name", language.Spanish, 1)
		if err != nil {
			fmt.Println(err)
		}
	}

	// reload caches
	buildCache(ctx)

	return []string{"Paises insertados/modificados. Cantidad: " + strconv.Itoa(num)}
}

func readFile(filename string) *xdominion.XRecords {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	data := &xdominion.XRecords{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		rec := scanLine(line)
		data.Push(rec)
	}
	return data
}

func scanLine(line string) *xdominion.XRecord {
	data := &xdominion.XRecord{}

	var fields map[string]interface{}
	err := json.Unmarshal([]byte(line), &fields)
	if err != nil {
		return nil
	}

	for i, v := range fields {
		data.Set(i, v)
	}
	return data
}
