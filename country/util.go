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

func buildTables(sitecontext *context.Context, databasename string) {

	for _, tbl := range moduletablesorder {
		sitecontext.Tables[tbl] = moduletables[tbl]()
		sitecontext.Tables[tbl].SetBase(sitecontext.Databases[databasename])
		sitecontext.Tables[tbl].SetLanguage(language.English)
	}
}

func createTables(sitecontext *context.Context) []string {

	messages := []string{}

	for _, tbl := range moduletablesorder {
		messages = append(messages, "Analysing "+tbl+" table.")
		num, err := sitecontext.Tables[tbl].Count(nil)
		if err != nil || num == 0 {
			err1 := sitecontext.Tables[tbl].Synchronize()
			if err1 != nil {
				messages = append(messages, "The table "+tbl+" was not created: "+err1.Error())
			} else {
				messages = append(messages, "The table "+tbl+" was created (again)")
			}
		} else {
			messages = append(messages, "The table "+tbl+" was not created because it contains data.")
		}
	}

	return messages
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
			key, _ := m.GetString("key")
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
		key, _ := r.(*xdominion.XRecord).GetString("key")
		sitecontext.Tables["country_country"].Upsert(key, *r.(*xdominion.XRecord))
		num++
	}

	// insert into translation each country in spanish
	DMPCOUNTRY = filespath + "countries.es.dmp"
	data = readFile(DMPCOUNTRY)
	for _, r := range *data {
		key, _ := r.(*xdominion.XRecord).GetString("key")
		name, _ := r.(*xdominion.XRecord).GetString("name")

		err := translation.SetTranslation(sitecontext, name, TRANSLATIONTHEME, key, "name", language.Spanish, 1)
		if err != nil {
			fmt.Println(err)
		}
	}

	// reload caches
	buildCache(sitecontext)

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
