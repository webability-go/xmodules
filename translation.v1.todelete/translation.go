package translation

import (
	gcontext "context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	//	"github.com/webability-go/xmodules/base"

	"github.com/webability-go/xamboo/applications"
)

const (
	SOURCETABLE = 1
	SOURCEFILE  = 2
)

func AddTheme(sitecontext applications.Datasource, theme string, name string, source int, link string, fields string) error {

	translation_theme := sitecontext.GetTable("translation_theme")
	if translation_theme == nil {
		sitecontext.Log("main", "xmodules::translation::addTheme: Error, the translation_theme table is not available on this context")
		return errors.New("xmodules::translation::addTheme: Error, the translation_theme table is not available on this context")
	}

	_, err := translation_theme.Upsert(theme, xdominion.XRecord{
		"key":    theme,
		"name":   name,
		"source": source,
		"link":   link,
		"fields": fields,
	})
	return err
}

// return: translation, ok (true, false), lastdate, lastverified (0, 1, 2)
// ok = true: texto correcto, false = no existe el texto
// last date = fecha en la cual se tradujo ( si no es español y ok = true)
// lastverified = 0: auto (o español original), 1 = verified, 2 = original modified (not re-translated, pending)
func GetTranslation(sitecontext applications.Datasource, textooriginal string, theme string, key string, field string, lang language.Tag) (string, bool, time.Time, int) {

	translation_info := sitecontext.GetTable("translation_info")
	if translation_info == nil {
		sitecontext.Log("main", "xmodules::translation::GetTranslation: Error, the translation_info table is not available on this context")
		return "", false, time.Time{}, 0
	}

	data, err := translation_info.SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("theme", "=", theme),
		xdominion.NewXCondition("externalkey", "=", key, "and"),
		xdominion.NewXCondition("field", "=", field, "and"),
		xdominion.NewXCondition("language", "=", lang.String(), "and"),
	})
	if err != nil {
		sitecontext.Log("main", "Error con select info", err)
		return "", false, time.Time{}, 0
	}

	if data != nil {
		lastdate, _ := data.GetTime("lastmodif")
		verify, _ := data.GetInt("verified")
		translation, _ := data.GetString("translation")
		return translation, true, lastdate, verify
	}

	return fmt.Sprintf("##%d::%s::%s##", theme, key, field), false, time.Time{}, 0
}

// return: error
func SetTranslation(sitecontext applications.Datasource, textotraducido string, theme string, key string, field string, lang language.Tag, verified int) error {

	translation_info := sitecontext.GetTable("translation_info")
	if translation_info == nil {
		sitecontext.Log("main", "xmodules::translation::SetTranslation: Error, the translation_info table is not available on this context")
		return errors.New("xmodules::translation::SetTranslation: Error, the translation_info table is not available on this context")
	}

	data, err := translation_info.SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("theme", "=", theme),
		xdominion.NewXCondition("externalkey", "=", key, "and"),
		xdominion.NewXCondition("field", "=", field, "and"),
		xdominion.NewXCondition("language", "=", lang.String(), "and"),
	})
	if err != nil {
		return err
	}

	if data != nil {
		// update
		key, _ := data.GetInt("key")
		_, err := translation_info.Update(key,
			xdominion.XRecord{
				"verify":      verified,
				"translation": textotraducido,
				"lastdate":    time.Now(),
				"lastuser":    1,
			})
		if err != nil {
			return err
		}
	} else {
		// insert
		_, err := translation_info.Insert(
			xdominion.XRecord{
				"key":         0,
				"theme":       theme,
				"language":    lang.String(),
				"externalkey": key,
				"field":       field,
				"translation": textotraducido,
				"lastmodif":   time.Now(),
				"lastuser":    1,
				"verified":    verified,
			})
		if err != nil {
			return err
		}
	}

	return nil
}

// return: error
func SetVerified(sitecontext applications.Datasource, theme string, key string, field string, lang language.Tag, verified int) error {

	translation_info := sitecontext.GetTable("translation_info")
	if translation_info == nil {
		sitecontext.Log("main", "xmodules::translation::SetVerified: Error, the translation_info table is not available on this context")
		return errors.New("xmodules::translation::SetVerified: Error, the translation_info table is not available on this context")
	}

	_, err := translation_info.Update(xdominion.XConditions{
		xdominion.NewXCondition("theme", "=", theme),
		xdominion.NewXCondition("externalkey", "=", key, "and"),
		xdominion.NewXCondition("field", "=", field, "and"),
		xdominion.NewXCondition("language", "=", lang.String(), "and"),
	},
		xdominion.XRecord{"verified": verified})
	return err
}

// Las credenciales de conección de google estan dentro del directorio accesible por GO en el archivo JSON de credenciales service_account
func GoogleTranslation(data []string, fromLang language.Tag, toLang language.Tag) ([]translate.Translation, error) {
	ctxbg := gcontext.Background()
	client, err := translate.NewClient(ctxbg)
	if err != nil {
		return nil, err
	}

	resp, err := client.Translate(ctxbg, data, toLang, &translate.Options{Source: fromLang})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
