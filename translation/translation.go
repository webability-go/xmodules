package translation

import (
	gcontext "context"
	"fmt"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"

	"xmodules/context"
)

func InitTranslation(sitecontext *context.Context, databasename string) error {

	sitecontext.Tables["kl_traducciontema"] = kl_traducciontema()
	sitecontext.Tables["kl_traducciontema"].SetBase(sitecontext.Databases[databasename])

	sitecontext.Tables["kl_traducciontabla"] = kl_traducciontabla()
	sitecontext.Tables["kl_traducciontabla"].SetBase(sitecontext.Databases[databasename])

	return nil
}

// return: translation, ok (true, false), lastdate, lastverified (0, 1, 2)
// ok = true: texto correcto, false = no existe el texto
// last date = fecha en la cual se tradujo ( si no es español y ok = true)
// lastverified = 0: auto (o español original), 1 = verified, 2 = original modified (not re-translated, pending)
func GetTraduccion(sitecontext *context.Context, textooriginal string, tema int, clave string, campo string, lang language.Tag) (string, bool, time.Time, int) {

	data, err := sitecontext.Tables["kl_traducciontabla"].SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("tema", "=", tema),
		xdominion.NewXCondition("claveext", "=", clave, "and"),
		xdominion.NewXCondition("campo", "=", campo, "and"),
		xdominion.NewXCondition("idioma", "=", lang.String(), "and"),
	})
	if err != nil {
		return "", false, time.Time{}, 0
	}

	if data != nil {
		lastdate, _ := data.GetTime("fecha")
		verify, _ := data.GetInt("verify")
		translation, _ := data.GetString("traduccion")
		return translation, true, lastdate, verify
	}

	return fmt.Sprintf("##%d::%s::%s##", tema, clave, campo), false, time.Time{}, 0
}

// return: error
func SetTraduccion(sitecontext *context.Context, textotraducido string, tema int, clave string, campo string, lang language.Tag, verified int) error {

	data, err := sitecontext.Tables["kl_traducciontabla"].SelectOne(xdominion.XConditions{
		xdominion.NewXCondition("tema", "=", tema),
		xdominion.NewXCondition("claveext", "=", clave, "and"),
		xdominion.NewXCondition("campo", "=", campo, "and"),
		xdominion.NewXCondition("idioma", "=", lang.String(), "and"),
	})
	if err != nil {
		return err
	}

	if data != nil {
		// update
		clave, _ := data.GetInt("clave")
		_, err := sitecontext.Tables["kl_traducciontabla"].Update(clave,
			xdominion.XRecord{
				"verify":     verified,
				"traduccion": textotraducido,
				"fecha":      time.Now(),
				"lastuser":   1,
			})
		if err != nil {
			return err
		}
	} else {
		// insert
		_, err := sitecontext.Tables["kl_traducciontabla"].Insert(
			xdominion.XRecord{
				"clave":      0,
				"tema":       tema,
				"idioma":     lang.String(),
				"claveext":   clave,
				"campo":      campo,
				"traduccion": textotraducido,
				"fecha":      time.Now(),
				"lastuser":   1,
				"verify":     verified,
			})
		if err != nil {
			return err
		}
	}

	return nil
}

// return: error
func SetVerified(sitecontext *context.Context, tema int, clave string, campo string, lang language.Tag, verified int) error {
	_, err := sitecontext.Tables["kl_traducciontabla"].Update(xdominion.XConditions{
		xdominion.NewXCondition("tema", "=", tema),
		xdominion.NewXCondition("claveext", "=", clave, "and"),
		xdominion.NewXCondition("campo", "=", campo, "and"),
		xdominion.NewXCondition("idioma", "=", lang.String(), "and"),
	},
		xdominion.XRecord{"verify": verified})
	return err
}

// Las credenciales de conección de google estan dentro del directorio accesible por GO en el archivo JSON de credenciales service_account
func GoogleTranslation(data []string, fromLang language.Tag, toLang language.Tag) ([]translate.Translation, error) {
	ctxbg := gcontext.Background()

	fmt.Println("Traduciendo: ", data, toLang)

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
