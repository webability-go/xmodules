package translation

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	//	"github.com/webability-go/xmodules/base"

	"github.com/webability-go/xamboo/applications"
)

type TranslationBlock struct {
	tema       int
	clave      string
	lastmodif  time.Time
	fromlang   language.Tag
	tolang     language.Tag
	prompt     string
	original   map[string]string
	verified   bool
	translated map[string]string
}

func NewTranslationBlock(tema int, clave string, lastmodif time.Time, fromlang language.Tag, tolang language.Tag) *TranslationBlock {
	return &TranslationBlock{
		tema:       tema,
		clave:      clave,
		lastmodif:  lastmodif,
		fromlang:   fromlang,
		tolang:     tolang,
		prompt:     "",
		original:   map[string]string{},
		verified:   false,
		translated: map[string]string{},
	}
}

func (tb *TranslationBlock) Set(field string, value string) {
	tb.original[field] = value
}

func (tb *TranslationBlock) SetPrompt(prompt string) {
	tb.prompt = prompt
}

func (tb *TranslationBlock) Verify(ds applications.Datasource) {

	data, err := ds.GetTable("kl_traducciontabla").SelectAll(xdominion.XConditions{
		xdominion.NewXCondition("tema", "=", tb.tema),
		xdominion.NewXCondition("claveext", "=", tb.clave, "and"),
		xdominion.NewXCondition("idioma", "=", tb.tolang.String(), "and"),
	})
	if err != nil {
		return
	}

	datafields := map[string]*xdominion.XRecord{}
	for _, rec := range *data {
		campo, _ := rec.GetString("campo")
		datafields[campo] = rec.(*xdominion.XRecord)
	}

	fields := []string{}
	values := []string{}

	for field, value := range tb.original {
		tr, ok := datafields[field]
		if !ok {
			fields = append(fields, field)
			values = append(values, value)
			continue
		}
		fecha, _ := tr.GetTime("fecha")
		if tb.lastmodif.After(fecha) {
			verify, _ := tr.GetInt("verify")
			if verify == 0 {
				fields = append(fields, field)
				values = append(values, value)
				continue
			} else if verify == 1 {
				// notificamos que hubo un cambio al usuario
				SetVerified(ds, tb.tema, tb.clave, field, tb.tolang, 2)
			}
		}
		trval, _ := tr.GetString("traduccion")
		tb.translated[field] = trval
	}
	if len(fields) > 0 {
		//		result, _ := GoogleTranslation(values, tb.fromlang, tb.tolang)

		langname := "inglés"
		langdata := GetLanguageByKey(ds, tb.tolang.String())
		//		fmt.Println("LANGUAGE TO TRANSLATE: ", tb.tolang.String(), langdata)
		if langdata != nil {
			langname, _ = langdata.GetString("name")
			//			fmt.Println("LANGUAGE NAME = ", langname)
		}
		if tb.prompt == "" {
			tb.prompt = "Traduce las líneas siguientes al {{LANGNAME}}, para una página web de recetas de cocina, guardando el mismo formato y sin quitar los números ni los | y sin poner un prompt en la respuesta:"
		}
		tb.prompt = strings.ReplaceAll(tb.prompt, "{{LANGNAME}}", langname)
		result, _ := GPTTranslation(ds, values, tb.fromlang, tb.prompt)

		if len(result) == len(fields) {
			for i, field := range fields {
				tb.translated[field] = result[i]
				//				tb.translated[field] = result[i].Text
				SetTraduccion(ds, result[i], tb.tema, tb.clave, field, tb.tolang, 0)
				//				SetTraduccion(sitecontext, result[i].Text, tb.tema, tb.clave, field, tb.tolang, 0)
			}
		} else {
			fmt.Println("TRADUCCION NO FUNCIONO: ", tb)
		}
	}
}

func (tb *TranslationBlock) Get(field string) string {
	return tb.translated[field]
}
