package translation

import (
	gcontext "context"
	"fmt"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

const (
	MODULEID = "translation"
	VERSION  = "1.0.0"

	SOURCETABLE = 1
	SOURCEFILE  = 2
)

func InitModule(sitecontext *context.Context, databasename string) error {

	buildTables(sitecontext, databasename)
	sitecontext.Modules[MODULEID] = VERSION

	return nil
}

func SynchronizeModule(sitecontext *context.Context) []string {

	messages := []string{}

	// Needed modules: context
	vc := context.ModuleInstalledVersion(sitecontext, "context")
	if vc == "" {
		messages = append(messages, "xmodules/context need to be installed before installing xmodules/translation.")
		return messages
	}

	messages = append(messages, "Analysing translation_theme table.")
	num, err := sitecontext.Tables["translation_theme"].Count(nil)
	if err != nil || num == 0 {
		err1 := sitecontext.Tables["translation_theme"].Synchronize()
		if err1 != nil {
			messages = append(messages, "The table translation_theme was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table translation_theme was created (again)")
		}
	} else {
		messages = append(messages, "The table translation_theme was not created because it contains data.")
	}

	messages = append(messages, "Analysing translation_info table.")
	num, err = sitecontext.Tables["translation_info"].Count(nil)
	if err != nil || num == 0 {
		err1 := sitecontext.Tables["translation_info"].Synchronize()
		if err1 != nil {
			messages = append(messages, "The table translation_info was not created: "+err1.Error())
		} else {
			messages = append(messages, "The table translation_info was created (again)")
		}
	} else {
		messages = append(messages, "The table translation_info was not created because it contains data.")
	}

	// Inserting into context-modules
	// Be sure context module is on db: fill context module (we should get this from xmodule.conf)
	err = context.AddModule(sitecontext, MODULEID, "Multilanguages translation tables for Xamboo", VERSION)
	if err == nil {
		messages = append(messages, "The entry "+MODULEID+" was modified successfully in the modules table.")
	} else {
		messages = append(messages, "Error modifying the entry "+MODULEID+" in the modules table: "+err.Error())
	}

	return messages
}

func AddTheme(sitecontext *context.Context, theme string, name string, source int, link string, fields string) error {
	_, err := sitecontext.Tables["translation_theme"].Upsert(theme, xdominion.XRecord{
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
func GetTranslation(sitecontext *context.Context, textooriginal string, theme string, key string, field string, lang language.Tag) (string, bool, time.Time, int) {

	data, err := sitecontext.Tables["translation_info"].SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("theme", "=", theme),
		xdominion.NewXCondition("externalkey", "=", key, "and"),
		xdominion.NewXCondition("field", "=", field, "and"),
		xdominion.NewXCondition("language", "=", lang.String(), "and"),
	})
	if err != nil {
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
func SetTranslation(sitecontext *context.Context, textotraducido string, theme string, key string, field string, lang language.Tag, verified int) error {

	data, err := sitecontext.Tables["translation_info"].SelectOne(xdominion.XConditions{
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
		_, err := sitecontext.Tables["translation_info"].Update(key,
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
		_, err := sitecontext.Tables["translation_info"].Insert(
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
func SetVerified(sitecontext *context.Context, theme string, key string, field string, lang language.Tag, verified int) error {
	_, err := sitecontext.Tables["translation_info"].Update(xdominion.XConditions{
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
