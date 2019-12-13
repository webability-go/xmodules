package translation

import (
	"time"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"

	"xmodules/context"
)

type TranslationBlock struct {
	tema       int
	clave      string
	lastmodif  time.Time
	fromlang   language.Tag
	tolang     language.Tag
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
		original:   map[string]string{},
		verified:   false,
		translated: map[string]string{},
	}
}

func (tb *TranslationBlock) Set(field string, value string) {
	tb.original[field] = value
}

func (tb *TranslationBlock) Verify(sitecontext *context.Context) {

	data, err := sitecontext.Tables["kl_traducciontabla"].SelectAll(xdominion.XConditions{
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
				SetVerified(sitecontext, tb.tema, tb.clave, field, tb.tolang, 2)
			}
		}
		trval, _ := tr.GetString("traduccion")
		tb.translated[field] = trval
	}
	if len(fields) > 0 {
		result, _ := GoogleTranslation(values, tb.fromlang, tb.tolang)
		for i, field := range fields {
			tb.translated[field] = result[i].Text
			SetTraduccion(sitecontext, result[i].Text, tb.tema, tb.clave, field, tb.tolang, 0)
		}
	}
}

func (tb *TranslationBlock) Get(field string) string {
	return tb.translated[field]
}
