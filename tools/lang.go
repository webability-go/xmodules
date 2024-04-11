package tools

import (
	"embed"
	"errors"
	"fmt"
	"log"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"

	"github.com/webability-go/xamboo/loggers"
)

// FUNCTIONS MADE TO BE USED BY BACK CODE WITH AN ARRAY OF TRANSLATION TABLES BUILD WTITH BuildMessagesFS

// If the first parameter is a datasource, will take the default language of the datasource.
// If the first parameter is a Language, will take the language of the parameter.
// If there is nothing found for the language, will take the english language, then the id itself

var Language = language.English

func Message(messages *map[language.Tag]*xcore.XLanguage, id string, params ...interface{}) string {

	var lang language.Tag
	ok := false
	if len(params) > 0 {
		lang, ok = params[0].(language.Tag)
		if ok {
			params = params[1:]
		} else {
			lang = Language
		}
	}
	if (*messages)[lang] == nil {
		lang = Language
	}
	if (*messages)[lang] == nil {
		lang = language.English
	}
	msg := (*messages)[lang].Get(id)
	if msg == "" && lang != language.English && (*messages)[language.English] != nil {
		msg = (*messages)[language.English].Get(id)
	}
	if msg == "" {
		msg = id
	}

	if len(params) > 0 {
		msg = fmt.Sprintf(msg, params...)
	}
	return msg
}

func ErrorMessage(messages *map[language.Tag]*xcore.XLanguage, id string, params ...interface{}) error {

	message := Message(messages, id, params...)
	return errors.New(message)
}

func MessageLog(log *log.Logger, messages *map[language.Tag]*xcore.XLanguage, id string, params ...interface{}) string {

	message := Message(messages, id, params...)
	log.Println(message)
	return message
}

func ErrorMessageLog(log *log.Logger, messages *map[language.Tag]*xcore.XLanguage, id string, params ...interface{}) error {

	message := Message(messages, id, params...)
	log.Println(message)
	return errors.New(message)
}

// FUNCTIONS MADE TO BE USED BY FRONT PAGES WITH AN ALREADY SELECTED LANGUAGE
func LogMessage(log *log.Logger, lang *xcore.XLanguage, id string, params ...interface{}) string {

	msg := lang.Get(id)
	if msg == "" {
		msg = id
	}

	if len(params) > 0 {
		msg = fmt.Sprintf(msg, params...)
	}
	log.Println(msg)
	return msg
}

func WajafLogErrorMessage(log *log.Logger, lang *xcore.XLanguage, id string, params ...interface{}) string {

	msg := LogMessage(log, lang, id, params...)
	return "{\"error\":true,\"message\":\"" + msg + "\"}"
}

func BuildMessagesFS(fs embed.FS, path string) *map[language.Tag]*xcore.XLanguage {

	// TODO(Phil) this is a workaround for the GetCoreLogger when the logger does not exists, it gives a panic error !!!
	slg := loggers.Create("errors", "stdout:", nil, nil).Logger
	if loggers.Loggers["X[errors]"] != nil {
		slg = loggers.GetCoreLogger("errors")
	}
	messages := &map[language.Tag]*xcore.XLanguage{}

	files, _ := fs.ReadDir(path)
	for _, file := range files {
		name := file.Name()
		pathname := name
		if path != "." {
			pathname = path + "/" + pathname
		}
		data, _ := fs.ReadFile(pathname)
		xlanguage, err := xcore.NewXLanguageFromXMLString(string(data))
		if err != nil {
			slg.Println("Error reading module messages:", pathname, err)
			continue
		}
		lang := xlanguage.GetLanguage()
		(*messages)[lang] = xlanguage
	}
	return messages
}
