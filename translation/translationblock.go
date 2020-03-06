package translation

import (
	"time"

	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

type TranslationBlock struct {
	tema       string
	clave      string
	lastmodif  time.Time
	fromlang   language.Tag
	tolang     language.Tag
	original   map[string]string
	verified   bool
	translated map[string]string
}

func NewTranslationBlock(tema string, clave string, lastmodif time.Time, fromlang language.Tag, tolang language.Tag) *TranslationBlock {
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

	translation_info := sitecontext.GetTable("translation_info")
	if translation_info == nil {
		sitecontext.Log("main", "xmodules::translation::Verify: Error, the translation_info table is not available on this context")
		return
	}

	data, err := translation_info.SelectAll(xdominion.XConditions{
		xdominion.NewXCondition("theme", "=", tb.tema),
		xdominion.NewXCondition("externalkey", "=", tb.clave, "and"),
		xdominion.NewXCondition("language", "=", tb.tolang.String(), "and"),
	})
	if err != nil {
		return
	}

	datafields := map[string]*xdominion.XRecord{}
	for _, rec := range *data {
		campo, _ := rec.GetString("field")
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
		fecha, _ := tr.GetTime("lastmodif")
		if tb.lastmodif.After(fecha) {
			verify, _ := tr.GetInt("verified")
			if verify == 0 {
				fields = append(fields, field)
				values = append(values, value)
				continue
			} else if verify == 1 {
				// notificamos que hubo un cambio al usuario
				SetVerified(sitecontext, tb.tema, tb.clave, field, tb.tolang, 2)
			}
		}
		trval, _ := tr.GetString("translation")
		tb.translated[field] = trval
	}
	if len(fields) > 0 {
		result, _ := GoogleTranslation(values, tb.fromlang, tb.tolang)
		for i, field := range fields {
			tb.translated[field] = result[i].Text
			SetTranslation(sitecontext, result[i].Text, tb.tema, tb.clave, field, tb.tolang, 0)
		}
	}
}

func (tb *TranslationBlock) Get(field string) string {
	return tb.translated[field]
}
