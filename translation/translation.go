package translation

import (
	gcontext "context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/chatgpt"

	"github.com/webability-go/xamboo/applications"
)

// TO REMOVE

func InitTranslation(sitecontext applications.Datasource, databasename string) error {

	sitecontext.SetTable("kl_traducciontema", kl_traducciontema())
	sitecontext.GetTable("kl_traducciontema").SetBase(sitecontext.GetDatabase())

	sitecontext.SetTable("kl_traducciontabla", kl_traducciontabla())
	sitecontext.GetTable("kl_traducciontabla").SetBase(sitecontext.GetDatabase())

	return nil
}

// ADMIN

func GetTranslationsCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	kl_traducciontabla := ds.GetTable("kl_traducciontabla")
	if kl_traducciontabla == nil {
		ds.Log("xmodules::video::GetCountProfilees: Error, the kl_traducciontabla table is not available on this datasource")
		return 0
	}
	cnt, _ := kl_traducciontabla.Count(cond)
	return cnt
}

func GetTranslationsList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	kl_traducciontabla := ds.GetTable("kl_traducciontabla")
	if kl_traducciontabla == nil {
		ds.Log("xmodules::video::GetProfileesList: Error, the kl_traducciontabla table is not available on this datasource")
		return nil
	}
	data, _ := kl_traducciontabla.SelectAll(cond, order, quantity, first)
	return data
}

func DeleteTranslationChildren(ds applications.Datasource, skey string) error {

	return nil
}

func PruneTranslationChildren(ds applications.Datasource, skey string, channel string) error {

	return nil
}

// return: translation, ok (true, false), lastdate, lastverified (0, 1, 2)
// ok = true: texto correcto, false = no existe el texto
// last date = fecha en la cual se tradujo ( si no es español y ok = true)
// lastverified = 0: auto (o español original), 1 = verified, 2 = original modified (not re-translated, pending)
func GetTraduccion(sitecontext applications.Datasource, textooriginal string, tema int, clave string, campo string, lang language.Tag) (string, bool, time.Time, int) {

	data, err := sitecontext.GetTable("kl_traducciontabla").SelectOne(xdominion.XConditions{
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
func SetTraduccion(sitecontext applications.Datasource, textotraducido string, tema int, clave string, campo string, lang language.Tag, verified int) error {

	data, err := sitecontext.GetTable("kl_traducciontabla").SelectOne(xdominion.XConditions{
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
		_, err := sitecontext.GetTable("kl_traducciontabla").Update(clave,
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
		_, err := sitecontext.GetTable("kl_traducciontabla").Insert(
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
func SetVerified(sitecontext applications.Datasource, tema int, clave string, campo string, lang language.Tag, verified int) error {
	_, err := sitecontext.GetTable("kl_traducciontabla").Update(xdominion.XConditions{
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

	//	fmt.Println("Traduciendo: ", data, toLang)

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

type choice struct {
	index   int `json:"index"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Finish_reason string `json:"finish_reason"`
}

type structresult struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
	Usage   struct {
		Prompt_tokens     int `json:"prompt_tokens"`
		Completion_tokens int `json:"completion_tokens"`
		Total_tokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Las credenciales de conección de google estan dentro del directorio accesible por GO en el archivo JSON de credenciales service_account
func GPTTranslation(ds applications.Datasource, data []string, fromLang language.Tag, prompt string) ([]string, error) {

	totranslate := ""
	for i, l := range data {
		// Condicionar la DATA:
		// retornos de linea codificados, pipe cambiado si hay pipe en la linea a traducir
		if strings.Contains(l, "\n") {
			l = strings.ReplaceAll(l, "\n", "<KBR>")
			l = strings.ReplaceAll(l, "\r", "")
		}
		if strings.Contains(l, "|") {
			l = strings.ReplaceAll(l, "|", "<KPP>")
		}
		totranslate += fmt.Sprint(i) + "|" + l + "\n"
	}
	//	fmt.Println("Traduciendo GPT: ", toLang, totranslate)

	result, err := chatgpt.TranslatePrompt(ds, prompt, totranslate)
	if err != nil {
		return nil, err
	}

	// 4. save translation into the file
	var sr structresult
	err = json.Unmarshal([]byte(result), &sr)
	if err != nil {
		fmt.Println("ERROR RESULTADO GPT: ", err)
		return nil, err
	}

	if len(sr.Choices) == 0 {
		fmt.Println("ERROR RESULTADO GPT: no hay choices en la respuesta")
		return nil, errors.New("GPT translation error: no hay choices en la respuesta")
	}

	messages := ""
	translated := []string{}
	aresult := strings.Split(sr.Choices[0].Message.Content, "\n")
	for _, l := range aresult {
		bresult := strings.Split(l, "|")
		if len(bresult) < 2 {
			messages += "Error result not correct: " + l + "\n"
			continue
		}
		if len(bresult) > 2 {
			bresult[1] = strings.Join(bresult[1:], "|")
		}
		i, _ := strconv.Atoi(bresult[0])
		if i >= len(data) || data[i] == "" {
			messages += "Error index result not correct: " + l + "\n"
			continue
		}
		nl := strings.ReplaceAll(bresult[1], "<KPP>", "|")
		nl = strings.ReplaceAll(nl, "<KBR>", "\n")
		translated = append(translated, nl)
	}
	fmt.Println("RESULTADO GPT: ", prompt, totranslate, sr.Choices[0].Message.Content, len(sr.Choices[0].Message.Content), " letras,", len(translated), " lineas.")

	return translated, nil
}
