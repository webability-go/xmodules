package country

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["country_country"] = countryCountry()
	sitecontext.Tables["country_country"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["country_country"].SetLanguage(language.English)
}

func buildCache(sitecontext *context.Context) {

	// Loads all data in XCache
	countries, _ := sitecontext.Tables["country_country"].SelectAll()
	if countries == nil {
		return
	}

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["country:countries:"+canonical] = xcore.NewXCache("country:countries:"+canonical, 0, 0)

		all := []string{}
		for _, m := range *countries {
			// creates structure on language
			str := CreateStructureCountryByData(sitecontext, m.Clone(), lang)
			key, _ := m.GetString("clave")
			all = append(all, key)
			sitecontext.Caches["country:countries:"+canonical].Set(key, str)
		}
		sitecontext.Caches["country:countries:"+canonical].Set("all", all)
	}
}

func loadTables(sitecontext *context.Context, filespath string) []string {

	// borra toda la data porque la vamos a insertar de nuevo (si se puede: FK bloquea)
	sitecontext.Tables["country_country"].Delete(nil)

	// 4 archivos de importación
	DMPCOUNTRY := filespath + "countries.en.dmp"

	num := 0
	data := readFile(DMPCOUNTRY)
	for _, r := range *data {

		sitecontext.Tables["country_country"].Upsert(*r.(*xdominion.XRecord))
		num++
	}

	// spanish in language

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
