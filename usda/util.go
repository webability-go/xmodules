package usda

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/language"

	"github.com/webability-go/xcore"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

func buildTables(sitecontext *context.Context, databasename string) {

	sitecontext.Tables["usda_group"] = usda_group()
	sitecontext.Tables["usda_group"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["usda_group"].SetLanguage(language.English)

	sitecontext.Tables["usda_food"] = usda_food()
	sitecontext.Tables["usda_food"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["usda_food"].SetLanguage(language.English)

	sitecontext.Tables["usda_nutrient"] = usda_nutrient()
	sitecontext.Tables["usda_nutrient"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["usda_nutrient"].SetLanguage(language.English)

	sitecontext.Tables["usda_foodnutrient"] = usda_foodnutrient()
	sitecontext.Tables["usda_foodnutrient"].SetBase(sitecontext.Databases[databasename])
	sitecontext.Tables["usda_foodnutrient"].SetLanguage(language.English)
}

func createCache(sitecontext *context.Context) []string {

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()
		sitecontext.Caches["usda:nutrients:"+canonical] = xcore.NewXCache("usda:nutrients:"+canonical, 0, 0)
	}
	return []string{}
}

func buildCache(sitecontext *context.Context) []string {

	// Loads all data in XCache
	nutrients, _ := sitecontext.Tables["usda_nutrient"].SelectAll()

	for _, lang := range sitecontext.Languages {
		canonical := lang.String()

		all := []string{}
		for _, m := range *nutrients {
			// creates structure on language
			str := CreateStructureNutrientByData(sitecontext, m.Clone(), lang)
			key, _ := m.GetString("key")
			all = append(all, key)
			sitecontext.Caches["usda:nutrients:"+canonical].Set(key, str)
		}
		sitecontext.Caches["usda:nutrients:"+canonical].Set("all", all)
	}

	return []string{}
}

func createTables(sitecontext *context.Context) []string {
	// alguna protección para saber si existe la tabla y no tronarla si tiene datos?
	// hacer un select count
	num1, err1 := sitecontext.Tables["usda_group"].Count(nil)
	num2, err2 := sitecontext.Tables["usda_food"].Count(nil)
	num3, err3 := sitecontext.Tables["usda_nutrient"].Count(nil)
	num4, err4 := sitecontext.Tables["usda_foodnutrient"].Count(nil)
	if (err1 != nil && err2 != nil && err3 != nil && err4 != nil) || (num1 == 0 && num2 == 0 && num3 == 0 && num4 == 0) {
		sitecontext.Logs["main"].Println("The tables usda_* were created (again)")
		sitecontext.Tables["usda_group"].Synchronize()
		sitecontext.Tables["usda_food"].Synchronize()
		sitecontext.Tables["usda_nutrient"].Synchronize()
		sitecontext.Tables["usda_foodnutrient"].Synchronize()
	} else {
		sitecontext.Logs["main"].Println("The tables usda_* were not created because they contain data")
	}

	return []string{}
}

func loadTables(sitecontext *context.Context, filespath string) []string {

	// borra toda la data porque la vamos a insertar de nuevo
	sitecontext.Tables["usda_foodnutrient"].Delete(nil)
	sitecontext.Tables["usda_nutrient"].Delete(nil)
	sitecontext.Tables["usda_food"].Delete(nil)
	sitecontext.Tables["usda_group"].Delete(nil)

	// 4 archivos de importación
	CSV_GROUP := filespath + "FD_GROUP.txt"
	CSV_NUTRIENT := filespath + "NUTR_DEF.txt"
	CSV_FOOD := filespath + "FOOD_DES.txt"
	CSV_FOODNUTRIENT := filespath + "NUT_DATA.txt"

	data := readFile(CSV_GROUP, map[int]string{
		0: "key",
		1: "name",
	})

	for _, r := range *data {
		r.Set("lastmodif", time.Now())

		sitecontext.Tables["usda_group"].Insert(*r.(*xdominion.XRecord))
	}
	// Adds group 9999
	sitecontext.Tables["usda_group"].Insert(xdominion.XRecord{
		"key":       "9999",
		"name":      "Other",
		"lastmodif": time.Now(),
	})

	data = readFile(CSV_NUTRIENT, map[int]string{
		0: "key",
		1: "unit",
		2: "tag",
		3: "name",
		4: "decimal",
		5: "order",
	})

	for _, r := range *data {
		r.Set("lastmodif", time.Now())
		// tag es opcional
		tag, _ := r.GetString("tag")
		if tag == "" {
			r.Set("tag", nil)
		}
		// decimal y order son int
		decimal, _ := r.GetString("decimal")
		i, _ := strconv.Atoi(decimal)
		r.Set("decimal", i)
		order, _ := r.GetString("order")
		o, _ := strconv.Atoi(order)
		r.Set("order", o)

		sitecontext.Tables["usda_nutrient"].Insert(*r.(*xdominion.XRecord))
	}
	// Adds nutrient 999: weight
	sitecontext.Tables["usda_nutrient"].Insert(xdominion.XRecord{
		"key":       "999",
		"name":      "Weight",
		"unit":      "g",
		"decimal":   1,
		"order":     1000,
		"lastmodif": time.Now(),
	})

	data = readFile(CSV_FOOD, map[int]string{
		0:  "key",
		1:  "group",
		2:  "name",
		3:  "abbr",
		4:  "commonname",
		10: "nfactor",
		11: "profactor",
		12: "fatfactor",
		13: "chofactor",
	})

	for _, r := range *data {
		r.Set("lastmodif", time.Now())
		// commonname es opcional
		tag, _ := r.GetString("commonname")
		if tag == "" {
			r.Set("commonname", nil)
		}
		// factores son  float64 o nil
		nfactor, _ := r.GetString("nfactor")
		if nfactor == "" {
			r.Set("nfactor", nil)
		} else {
			f, _ := strconv.ParseFloat(nfactor, 64)
			r.Set("nfactor", f)
		}
		profactor, _ := r.GetString("profactor")
		if profactor == "" {
			r.Set("profactor", nil)
		} else {
			f, _ := strconv.ParseFloat(profactor, 64)
			r.Set("profactor", f)
		}
		fatfactor, _ := r.GetString("fatfactor")
		if fatfactor == "" {
			r.Set("fatfactor", nil)
		} else {
			f, _ := strconv.ParseFloat(fatfactor, 64)
			r.Set("fatfactor", f)
		}
		chofactor, _ := r.GetString("chofactor")
		if chofactor == "" {
			r.Set("chofactor", nil)
		} else {
			f, _ := strconv.ParseFloat(chofactor, 64)
			r.Set("chofactor", f)
		}

		sitecontext.Tables["usda_food"].Insert(*r.(*xdominion.XRecord))
	}

	data = readFile(CSV_FOODNUTRIENT, map[int]string{
		0: "food",
		1: "nutrient",
		2: "value",
	})

	for _, r := range *data {
		r.Set("key", 0)
		r.Set("lastmodif", time.Now())
		value, _ := r.GetString("value")
		if value == "" {
			r.Set("value", nil)
		} else {
			v, _ := strconv.ParseFloat(value, 64)
			r.Set("value", v)
		}

		sitecontext.Tables["usda_foodnutrient"].Insert(*r.(*xdominion.XRecord))
	}

	return []string{}
}

func readFile(filename string, fields map[int]string) *xdominion.XRecords {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	data := &xdominion.XRecords{}
	utf := charmap.ISO8859_1.NewDecoder().Reader(file)
	scanner := bufio.NewScanner(utf)
	for scanner.Scan() {
		line := scanner.Text()
		rec := scanLine(line, fields)
		data.Push(rec)
	}
	return data
}

func scanLine(line string, fields map[int]string) *xdominion.XRecord {
	data := &xdominion.XRecord{}

	xline := strings.Split(line, "^")
	for i, v := range fields {
		val := strings.Replace(xline[i], "~", "", -1)
		data.Set(v, val)
	}
	return data
}
