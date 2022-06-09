package tools

import (
	"embed"
	"fmt"
	"log"

	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"

	"github.com/webability-go/xamboo/loggers"
)

// Language is the default system language to search messages in this language in priority
// It may be changed by the code at anytime. Recommended to change it to your language when you init your libraries
var Language language.Tag = language.English

func Message(messages *map[language.Tag]*xcore.XLanguage, id string, params ...interface{}) string {

	lang := Language
	ok := false
	if len(params) > 0 {
		lang, ok = params[0].(language.Tag)
		if ok {
			params = params[1:]
		} else {
			lang = Language
		}
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

func BuildMessages(data map[language.Tag]map[string]string) *map[language.Tag]*xcore.XLanguage {
	bdata := map[language.Tag]*xcore.XLanguage{}
	for l, t := range data {
		xl := xcore.NewXLanguage("messages", l)
		for id, val := range t {
			xl.Set(id, val)
		}
		bdata[l] = xl
	}
	return &bdata
}

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

	slg := loggers.GetCoreLogger("errors")
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
			slg.Println("Error reading module messages:", err)
			continue
		}
		lang := xlanguage.GetLanguage()
		(*messages)[lang] = xlanguage
	}
	return messages
}
